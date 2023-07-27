/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package unsharded

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"vitess.io/vitess/go/test/endtoend/utils"
	"vitess.io/vitess/go/vt/sidecardb"

	"github.com/stretchr/testify/require"

	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/test/endtoend/cluster"
)

var (
	clusterInstance *cluster.LocalProcessCluster
	vtParamsKs1     mysql.ConnParams
	vtParamsKs2     mysql.ConnParams
	keyspace1Name   = "ks"
	keyspace2Name   = "ks2"
	sidecarDBName   = "_vt_schema_tracker_metadata" // custom sidecar database name for testing
	cell            = "zone1"
	sqlSchema       = `
		create table main (
			id bigint,
			val varchar(128),
			primary key(id)
		) Engine=InnoDB;
`
)

func TestMain(m *testing.M) {
	defer cluster.PanicHandler(nil)
	flag.Parse()

	exitCode := func() int {
		clusterInstance = cluster.NewCluster(cell, "localhost")
		defer clusterInstance.Teardown()

		vtgateVer, err := cluster.GetMajorVersion("vtgate")
		if err != nil {
			return 1
		}
		vttabletVer, err := cluster.GetMajorVersion("vttablet")
		if err != nil {
			return 1
		}

		// For upgrade/downgrade tests.
		if vtgateVer < 17 || vttabletVer < 17 {
			// Then only the default sidecarDBName is supported.
			sidecarDBName = sidecardb.DefaultName
		}

		// Start topo server
		err = clusterInstance.StartTopo()
		if err != nil {
			return 1
		}

		clusterInstance.VtTabletExtraArgs = []string{"--queryserver-config-schema-change-signal", "--watch_replication_stream", "--track_schema_versions"}

		// Start keyspace1
		keyspace1 := &cluster.Keyspace{
			Name:          keyspace1Name,
			SchemaSQL:     sqlSchema,
			SidecarDBName: sidecarDBName,
		}
		err = clusterInstance.StartUnshardedKeyspace(*keyspace1, 0, false)
		if err != nil {
			return 1
		}

		// Start keyspace2
		keyspace2 := &cluster.Keyspace{
			Name:          keyspace2Name,
			SchemaSQL:     sqlSchema,
			SidecarDBName: sidecarDBName,
		}
		err = clusterInstance.StartUnshardedKeyspace(*keyspace2, 0, false)
		if err != nil {
			return 1
		}

		// Start vtgate
		clusterInstance.VtGateExtraArgs = []string{"--schema_change_signal", "--vschema_ddl_authorized_users", "%"}
		err = clusterInstance.StartVtgate()
		if err != nil {
			return 1
		}

		err = clusterInstance.WaitForVTGateAndVTTablets(5 * time.Minute)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		vtParamsKs1 = mysql.ConnParams{
			Host:   clusterInstance.Hostname,
			Port:   clusterInstance.VtgateMySQLPort,
			DbName: keyspace1Name,
		}

		vtParamsKs2 = mysql.ConnParams{
			Host:   clusterInstance.Hostname,
			Port:   clusterInstance.VtgateMySQLPort,
			DbName: keyspace2Name,
		}

		return m.Run()
	}()
	os.Exit(exitCode)
}

func TestNewUnshardedTable(t *testing.T) {
	defer cluster.PanicHandler(t)

	// create a sql connection
	ctx := context.Background()
	connKs1, err := mysql.Connect(ctx, &vtParamsKs1)
	require.NoError(t, err)
	defer connKs1.Close()

	connKs2, err := mysql.Connect(ctx, &vtParamsKs2)
	require.NoError(t, err)
	defer connKs2.Close()

	vtgateVersion, err := cluster.GetMajorVersion("vtgate")
	require.NoError(t, err)
	expected := `[[VARCHAR("dual")] [VARCHAR("main")]]`
	if vtgateVersion >= 17 {
		expected = `[[VARCHAR("main")]]`
	}

	// ensuring our initial table "main" is in the schema
	utils.AssertMatchesWithTimeout(t, connKs1,
		"SHOW VSCHEMA TABLES",
		expected,
		100*time.Millisecond,
		30*time.Second,
		"initial table list not complete")
	utils.AssertMatchesWithTimeout(t, connKs2,
		"SHOW VSCHEMA TABLES",
		expected,
		100*time.Millisecond,
		30*time.Second,
		"initial table list not complete")

	// create a new table which is not part of the VSchema
	wg := sync.WaitGroup{}
	wg.Add(2)
	tableCount := 100

	go func() {
		defer wg.Done()
		for i := 0; i < tableCount; i++ {
			ddl := fmt.Sprintf("create table new_table_tracked1_%d(id bigint, name varchar(100), primary key(id)) Engine=InnoDB", i)
			utils.Exec(t, connKs1, ddl)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < tableCount; i++ {
			ddl := fmt.Sprintf("create table new_table_tracked2_%d(id bigint, name varchar(100), primary key(id)) Engine=InnoDB", i)
			utils.Exec(t, connKs2, ddl)
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()

	expected1 := `[[VARCHAR("main")]`
	for i := 0; i < tableCount; i++ {
		expected1 = expected1 + fmt.Sprintf(` [VARCHAR("new_table_tracked1_%d")]`, i)
	}
	expected1 = expected1 + `]`

	expected2 := `[[VARCHAR("main")]`
	for i := 0; i < tableCount; i++ {
		expected2 = expected2 + fmt.Sprintf(` [VARCHAR("new_table_tracked2_%d")]`, i)
	}
	expected2 = expected2 + `]`

	// waiting for the vttablet's schema_reload interval to kick in
	utils.AssertMatchesWithTimeout(t, connKs1,
		"SHOW VSCHEMA TABLES",
		expected1,
		100*time.Millisecond,
		30*time.Second,
		"new_table_tracked not in vschema tables")
	utils.AssertMatchesWithTimeout(t, connKs2,
		"SHOW VSCHEMA TABLES",
		expected2,
		100*time.Millisecond,
		30*time.Second,
		"new_table_tracked not in vschema tables")
}
