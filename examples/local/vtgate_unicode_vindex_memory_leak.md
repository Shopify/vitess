Setup
======
* Edit `go/vt/vtgate/vindexes/unicode.go` to not use a pool (just a single shared `pooledCollator` protected by a mutex)
  This ensures that reproduction isn't impacted by the nondeterministic nature of sync.Pool, i.e. https://pkg.go.dev/sync#Pool
  > Any item stored in the Pool may be removed automatically at any time without notification. If the Pool holds the only reference when this happens, the item might be deallocated.
  
* `SKIP_VTADMIN=true ./101_initial_cluster.sh`
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
            "unicode_loose_md5": {
                "type": "unicode_loose_md5"
            }
        },
        "tables": {
            "things": {
                "column_vindexes": [
                    {
                        "column": "id",
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
  create table if not exists things(
    id varchar(255) not null,
    primary key(id)
  ) ENGINE=InnoDB;
' things || fail "Failed to create tables in the things keyspace"
```

* Check vtgate memory use
```
cat /proc/$(cat /tmp/vtdataroot/tmp/vtgate.pid)/status | grep RSS
```

* Run lots of queries
```bash
go run run_queries.go
```

* Keep checking memory use
