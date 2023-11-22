package main

import (
	"context"
	"log"

	_ "vitess.io/vitess/go/vt/vtctl/grpcvtctlclient"
	_ "vitess.io/vitess/go/vt/vtgate/grpcvtgateconn"
	"vitess.io/vitess/go/vt/vtgate/vtgateconn"
)

func main() {
	ctx := context.Background()
	conn, err := vtgateconn.Dial(ctx, "localhost:15991")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	session := conn.Session("things", nil)
	for {
		_, err := session.Execute(ctx, "SELECT * FROM things WHERE id = 'fooBARfooBARfooBARfooBARfooBAR'", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
