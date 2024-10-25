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

package zk2topo

import (
	"context"
	"testing"

	"github.com/z-division/go-zookeeper/zk"
	"vitess.io/vitess/go/testfiles"
	"vitess.io/vitess/go/vt/zkctl"
)

func TestZkConnClosedOnDisconnect(t *testing.T) {
	zkd, serverAddr := zkctl.StartLocalZk(testfiles.GoVtTopoZk2topoZkID, testfiles.GoVtTopoZk2topoPort)
	defer zkd.Teardown()

	conn := Connect(serverAddr)
	_, _, err := conn.Get(context.Background(), "/")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if !conn.conn.State().IsConnected() {
		t.Fatalf("Connection not connected: %v", conn.conn.State())
	}

	oldConn := conn.conn

	// simulate a disconnect
	zkd.Shutdown()
	zkd.Start()

	// do another get to trigger a new connection
	_, _, err = conn.Get(context.Background(), "/")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// Check that old connection is closed
	_, _, err = oldConn.Get("/")
	if err == nil {
		t.Fatalf("Get() should have failed: %v", err)
	}

	if oldConn.State() != zk.StateDisconnected {
		t.Fatalf("Connection not closed: %v", oldConn.State())
	}
}
