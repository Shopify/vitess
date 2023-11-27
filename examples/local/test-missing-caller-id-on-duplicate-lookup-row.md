Setup
======
* `SKIP_VTADMIN=true ./101_initial_cluster.sh`
* Create a new sharded "things" keyspace:
```bash
source env.sh

for i in 300 301 302; do
	CELL=zone1 TABLET_UID=$i ./scripts/mysqlctl-up.sh
	SHARD=-80 CELL=zone1 KEYSPACE=things TABLET_UID=$i ./scripts/vttablet-up.sh
done

for i in 400 401 402; do
	CELL=zone1 TABLET_UID=$i ./scripts/mysqlctl-up.sh
	SHARD=80- CELL=zone1 KEYSPACE=things TABLET_UID=$i ./scripts/vttablet-up.sh
done

# set the correct durability policy for the keyspace
vtctldclient --server localhost:15999 SetKeyspaceDurabilityPolicy --durability-policy=semi_sync things || fail "Failed to set keyspace durability policy on the things keyspace"

for shard in "-80" "80-"; do
	# Wait for all the tablets to be up and registered in the topology server
	# and for a primary tablet to be elected in the shard and become healthy/serving.
	wait_for_healthy_shard things "${shard}" || exit 1
done
```
* Create the VSchema:
```bash
vtctldclient ApplyVSchema --vschema '
    {
        "sharded": true,
        "vindexes": {
            "hash": {
                "type": "hash"
            },
            "unicode_loose_md5": {
                "type": "unicode_loose_md5"
            },
            "things_name_lookup": {
                "type": "consistent_lookup_unique",
                "params": {
                    "table": "things_name_lookup",
                    "from": "name",
                    "to": "keyspace_id"
                },
                "owner": "things"
            }
        },
        "tables": {
            "things": {
                "column_vindexes": [
                    {
                        "column": "id",
                        "name": "hash"
                    },
                    {
                        "column": "name",
                        "name": "things_name_lookup"
                    }
                ]
            },
            "things_name_lookup": {
                "column_vindexes": [
                    {
                        "column": "name",
                        "name": "unicode_loose_md5"
                    }
                ]
            }            
        }
    }
    ' things || fail "Failed to create vschema in sharded things keyspace"
```

* Create the SQL schema:
```bash
vtctldclient ApplySchema --sql '
  CREATE TABLE things (
    id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
  ) ENGINE=InnoDB;

  CREATE TABLE things_name_lookup (
    name VARCHAR(255) NOT NULL,
    keyspace_id VARBINARY(10) NOT NULL,
    PRIMARY KEY (name)
  ) ENGINE=InnoDB;
' things || fail "Failed to create tables in the things keyspace"
```

* Log in to vtgate with password `mysql_password`
```
mysql --user mysql_user --password
```
* Insert some data
```
USE things;
INSERT INTO things (id, name) VALUES (1, "foo"), (2, "bar");

select * from things;
+----+------+
| id | name |
+----+------+
|  2 | bar  |
|  1 | foo  |
+----+------+
2 rows in set (0.00 sec)

select * from things_name_lookup;
+------+--------------------------+
| name | keyspace_id              |
+------+--------------------------+
| bar  | 0x06E7EA22CE92708F       |
| foo  | 0x166B40B44ABA4BD6       |
+------+--------------------------+
2 rows in set (0.00 sec)

INSERT INTO things (id, name) VALUES (3, "bar");
ERROR 1045 (28000): transaction rolled back to reverse changes of partial DML execution: target: things.-80.primary: vttablet: missing caller id
```

* Update `go/vt/vtgate/vindexes/consistent_lookup.go` line 364 to pass `ctx`
* Rebuild, restart vtgate, and try again
```
mysql> INSERT INTO things (id, name) VALUES (3, "bar");
ERROR 1062 (23000): transaction rolled back to reverse changes of partial DML execution: lookup.Create: Code: ALREADY_EXISTS
vttablet: Duplicate entry 'bar' for key 'things_name_lookup.PRIMARY' (errno 1062) (sqlstate 23000) (CallerID: mysql_user): Sql: "insert into things_name_lookup(`name`, keyspace_id) values (:_name_0, :keyspace_id_0)", BindVars: {_name_0: "type:VARCHAR value:\"bar\""keyspace_id_0: "type:VARBINARY value:\"N\\xb1\\x90ɢ\\xfa\\x16\\x9c\""name_0: "type:VARCHAR value:\"bar\""}
```
