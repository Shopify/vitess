#!/bin/bash

# Copyright 2019 The Vitess Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# this script brings up topo server and all the vitess components
# required for a single shard deployment.

source ./env.sh

# start topo server
if [ "${TOPO}" = "zk2" ]; then
	CELL=zone1 ./scripts/zk-up.sh
elif [ "${TOPO}" = "k8s" ]; then
	CELL=zone1 ./scripts/k3s-up.sh
else
	CELL=zone1 ./scripts/etcd-up.sh
fi

# start vtctld
CELL=zone1 ./scripts/vtctld-up.sh

# start vttablets for new shards. we start only one tablet each (primary)
CELL=zone1 TABLET_UID=200 ./scripts/mysqlctl-up.sh
SHARD=-40 CELL=zone1 KEYSPACE=main TABLET_UID=200 ./scripts/vttablet-up.sh
CELL=zone1 TABLET_UID=300 ./scripts/mysqlctl-up.sh
SHARD=40-80 CELL=zone1 KEYSPACE=main TABLET_UID=300 ./scripts/vttablet-up.sh
CELL=zone1 TABLET_UID=400 ./scripts/mysqlctl-up.sh
SHARD=80-c0 CELL=zone1 KEYSPACE=main TABLET_UID=400 ./scripts/vttablet-up.sh
CELL=zone1 TABLET_UID=500 ./scripts/mysqlctl-up.sh
SHARD=c0- CELL=zone1 KEYSPACE=main TABLET_UID=500 ./scripts/vttablet-up.sh

# set the correct durability policy for the keyspace
vtctldclient --server localhost:15999 SetKeyspaceDurabilityPolicy --durability-policy=none main

# start vtorc
./scripts/vtorc-up.sh

for shard in "-40" "40-80" "80-c0" "c0-"; do
  # Wait for a primary tablet to be elected in the shard
  for _ in $(seq 0 200); do
  	vtctldclient GetTablets --keyspace main --shard $shard | grep -q "primary" && break
  	sleep 1
  done;
  vtctldclient GetTablets --keyspace main --shard $shard | grep "primary" || (echo "Timed out waiting for primary to be elected in main/$shard" && exit 1)
done;


# create the vschema
vtctldclient ApplyVSchema --vschema-file main_vschema.json main

# create the schema
vtctldclient ApplySchema --sql-file create_main_schema.sql main

# start vtgate
CELL=zone1 ./scripts/vtgate-up.sh

# start vtadmin
#./scripts/vtadmin-up.sh

