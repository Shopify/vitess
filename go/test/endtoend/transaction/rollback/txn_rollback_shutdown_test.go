/*
Copyright 2019 The Vitess Authors.

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

package rollback

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"vitess.io/vitess/go/mysql/sqlerror"
	"vitess.io/vitess/go/test/endtoend/utils"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/test/endtoend/cluster"
)

var (
	clusterInstance *cluster.LocalProcessCluster
	vtParams        mysql.ConnParams
	keyspaceName    = "ks"
	cell            = "zone1"
	hostname        = "localhost"
	sqlSchema       = `
	create table buffer(
		id BIGINT NOT NULL,
		msg VARCHAR(64) NOT NULL,
		PRIMARY KEY (id)
	) Engine=InnoDB;`
)

func TestMain(m *testing.M) {
	defer cluster.PanicHandler(nil)
	flag.Parse()

	exitCode := func() int {
		clusterInstance = cluster.NewCluster(cell, hostname)
		defer clusterInstance.Teardown()

		// Reserve vtGate port in order to pass it to vtTablet
		clusterInstance.VtgateGrpcPort = clusterInstance.GetAndReservePort()

		// Start topo server
		err := clusterInstance.StartTopo()
		if err != nil {
			panic(err)
		}

		// Start keyspace
		keyspace := &cluster.Keyspace{
			Name:      keyspaceName,
			SchemaSQL: sqlSchema,
		}
		err = clusterInstance.StartUnshardedKeyspace(*keyspace, 1, false)
		if err != nil {
			panic(err)
		}

		// Set a short onterm timeout so the test goes faster.
		clusterInstance.VtGateExtraArgs = []string{"--onterm_timeout", "1s"}
		err = clusterInstance.StartVtgate()
		if err != nil {
			panic(err)
		}
		vtParams = clusterInstance.GetVTParams(keyspaceName)
		return m.Run()
	}()
	os.Exit(exitCode)
}

func TestTransactionRollBackWhenShutDown(t *testing.T) {
	defer cluster.PanicHandler(t)
	ctx := context.Background()
	conn, err := mysql.Connect(ctx, &vtParams)
	require.NoError(t, err)
	defer conn.Close()

	utils.Exec(t, conn, "insert into buffer(id, msg) values(3,'mark')")
	utils.Exec(t, conn, "insert into buffer(id, msg) values(4,'doug')")

	// start an incomplete transaction
	utils.Exec(t, conn, "begin")
	utils.Exec(t, conn, "select * from buffer where id = 3 for update")

	// Enforce a restart to enforce rollback
	if err = clusterInstance.RestartVtgate(); err != nil {
		t.Errorf("Fail to re-start vtgate: %v", err)
	}

	// Make a new mysql connection to vtGate
	vtParams = clusterInstance.GetVTParams(keyspaceName)
	conn2, err := mysql.Connect(ctx, &vtParams)
	require.NoError(t, err)
	defer conn2.Close()

	// Start a new transaction
	utils.Exec(t, conn2, "begin")
	defer utils.Exec(t, conn2, "rollback")
	// Verify previous transaction was rolled back. Row lock should be available, otherwise we'll get an error.
	qr := utils.Exec(t, conn2, "select * from buffer where id = 3 for update nowait")
	assert.Equal(t, 1, len(qr.Rows))

}

func TestTransactionRollBackWhenShutDownWithQueryRunning(t *testing.T) {
	defer cluster.PanicHandler(t)
	ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, err := mysql.Connect(ctx, &vtParams)
		require.NoError(t, err)
		defer conn.Close()

		utils.Exec(t, conn, "insert into buffer(id, msg) values(5,'alpha')")
		utils.Exec(t, conn, "insert into buffer(id, msg) values(6,'beta')")

		// start an incomplete transaction with a long-running query
		utils.Exec(t, conn, "begin")
		_, err = conn.ExecuteFetch("select *, sleep(40) from buffer where id = 5 for update", 1, true)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "EOF")
	}()

	// wait for the long-running query to start executing
	checkConn, err := mysql.Connect(ctx, &vtParams)
	require.NoError(t, err)
	defer checkConn.Close()
	waitTimeout := time.After(10 * time.Second)
	running := false
	for running == false {
		select {
		case <-waitTimeout:
			t.Fatalf("Long-running query did not start executing")
		case <-time.After(10 * time.Millisecond):
			// We should get a lock wait timeout error once the long-running query starts executing
			_, err := checkConn.ExecuteFetch("select * from buffer where id = 5 for update nowait", 1, true)
			if sqlErr, ok := err.(*sqlerror.SQLError); ok {
				if sqlErr.Number() == sqlerror.ERLockNowait {
					running = true
					continue
				}
			}
			require.NoError(t, err)
		}
	}

	// Enforce a restart to enforce rollback
	if err = clusterInstance.RestartVtgate(); err != nil {
		t.Errorf("Fail to re-start vtgate: %v", err)
	}

	// Make a new mysql connection to vtGate
	vtParams = clusterInstance.GetVTParams(keyspaceName)
	conn2, err := mysql.Connect(ctx, &vtParams)
	require.NoError(t, err)
	defer conn2.Close()

	// Verify previous transaction was rolled back. Row lock should be available, otherwise we'll get an error.
	qr := utils.Exec(t, conn2, "select * from buffer where id = 5 for update nowait")
	assert.Equal(t, 1, len(qr.Rows))
}

func TestErrorInAutocommitSession(t *testing.T) {
	defer cluster.PanicHandler(t)
	ctx := context.Background()
	conn, err := mysql.Connect(ctx, &vtParams)
	require.NoError(t, err)
	defer conn.Close()

	utils.Exec(t, conn, "set autocommit=true")
	utils.Exec(t, conn, "insert into buffer(id, msg) values(1,'foo')")
	_, err = conn.ExecuteFetch("insert into buffer(id, msg) values(1,'bar')", 1, true)
	require.Error(t, err) // this should fail with duplicate error
	utils.Exec(t, conn, "insert into buffer(id, msg) values(2,'baz')")

	conn2, err := mysql.Connect(ctx, &vtParams)
	require.NoError(t, err)
	defer conn2.Close()
	result := utils.Exec(t, conn2, "select * from buffer where id in (1,2) order by id")

	// if we have properly working autocommit code, both the successful inserts should be visible to a second
	// connection, even if we have not done an explicit commit
	assert.Equal(t, `[[INT64(1) VARCHAR("foo")] [INT64(2) VARCHAR("baz")]]`, fmt.Sprintf("%v", result.Rows))
}
