In examples/local
==================
* Create a `table_acl.json` file
```bash
echo '{
  "table_groups": [
    {
      "name": "all-tables",
      "table_names_or_prefixes": [
        "%"
      ],
      "admins": [
        "mysql_user"
      ],
      "readers": [
        "mysql_user"
      ],
      "writers": [
        "mysql_user"
      ]
    }
  ]
}' > table_acl.json
```
* Add authentication flags to `examples/common/scripts/vtgate-up.sh`
```
  --mysql_auth_server_impl static \
  --mysql_auth_server_static_file "$(dirname "${BASH_SOURCE[0]:-$0}")/../../local/mysql_auth_server_static_creds.json" \
  --mysql_auth_static_reload_interval 5s \
```
* Add strict ACL enforcement flags to `examples/common/scripts/vttablet-up.sh`
```
 --enforce-tableacl-config \
 --table-acl-config-reload-interval 5s \
 --queryserver-config-strict-table-acl \
 --table-acl-config "$(dirname "${BASH_SOURCE[0]:-$0}")/../../local/table_acl.json" \
```
* Start the cluster
```bash
SKIP_VTADMIN=true ./101_initial_cluster.sh
```
* Create a new sharded "things" keyspace:
```bash
source ../common/env.sh

for i in 300 301 302; do
	CELL=zone1 TABLET_UID=$i ../common/scripts/mysqlctl-up.sh
	SHARD=-80 CELL=zone1 KEYSPACE=things TABLET_UID=$i ../common/scripts/vttablet-up.sh
done

for i in 400 401 402; do
	CELL=zone1 TABLET_UID=$i ../common/scripts/mysqlctl-up.sh
	SHARD=80- CELL=zone1 KEYSPACE=things TABLET_UID=$i ../common/scripts/vttablet-up.sh
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
            "xxhash": {
                "type": "xxhash"
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
                        "column": "uuid",
                        "name": "xxhash"
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
    uuid VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (uuid),
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
mysql> USE things;

mysql> INSERT INTO things (uuid, name) VALUES (UUID(), "foo"), (UUID(), "bar");

mysql> SELECT * FROM things;
+--------------------------------------+------+
| uuid                                 | name |
+--------------------------------------+------+
| c9beff6d-8fe4-11ee-9f69-4201f31b044a | bar  |
| c9beff2a-8fe4-11ee-9f69-4201f31b044a | foo  |
+--------------------------------------+------+

mysql> SELECT * FROM things_name_lookup;
+------+--------------------------+
| name | keyspace_id              |
+------+--------------------------+
| foo  | 0xFEBD7FC7483D91F0       |
| bar  | 0x3043FD669C06B406       |
+------+--------------------------+
2 rows in set (0.00 sec)
```
* Try to insert a duplicate row
```
mysql> INSERT INTO things (uuid, name) VALUES (UUID(), "bar");
ERROR 1045 (28000): transaction rolled back to reverse changes of partial DML execution: target: things.-80.primary: vttablet: missing caller id
```
* Update `go/vt/vtgate/vindexes/consistent_lookup.go` `handleDup` function to pass `ctx` instead of `context.Background()`
* Rebuild, restart vtgate, and try again
```
mysql> INSERT INTO things (uuid, name) VALUES (UUID(), "bar");
ERROR 1062 (23000): transaction rolled back to reverse changes of partial DML execution: lookup.Create: Code: ALREADY_EXISTS
vttablet: Duplicate entry 'bar' for key 'things_name_lookup.PRIMARY' (errno 1062) (sqlstate 23000) (CallerID: mysql_user): Sql: "insert into things_name_lookup(`name`, keyspace_id) values (:_name_0, :keyspace_id_0)", BindVars: {_name_0: "type:VARCHAR value:\"bar\""keyspace_id_0: "type:VARBINARY value:\"\\x89W\\x9e\\xb4\\xdc\\xfb\\n\\xf0\""}

target: things.80-.primary
```
