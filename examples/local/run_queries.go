package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "vitess.io/vitess/go/vt/vtctl/grpcvtctlclient"
	_ "vitess.io/vitess/go/vt/vtgate/grpcvtgateconn"
	"vitess.io/vitess/go/vt/vtgate/vtgateconn"
)

const NumQueries = 10000
const Concurrency = 100
const ShouldDelay = false
const Delay = 100 * time.Millisecond

func main() {
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer connectCancel()

	fmt.Println("Connecting")
	conn, err := vtgateconn.Dial(connectCtx, "localhost:15991")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session := conn.Session("commerce", nil)

	var wg sync.WaitGroup
	wg.Add(Concurrency)

	tasks := make(chan int, NumQueries)

	for i := 0; i < Concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for range tasks {
				duration, err := runQuery(session)
				if err != nil {
					fmt.Printf("%v %v\n", duration, err)
					os.Exit(1)
				}
				//fmt.Println(duration)
				if ShouldDelay {
					time.Sleep(Delay)
				}
			}
		}(i)
	}

	// Enqueue X tasks
	for i := 0; i < NumQueries; i++ {
		tasks <- i
	}

	// Close the channel to signal to the goroutines that no more tasks will be added
	close(tasks)

	// Wait for all goroutines to finish
	wg.Wait()
}

func runQuery(session *vtgateconn.VTGateSession) (time.Duration, error) {
	queryCtx, queryCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer queryCancel()
	start := time.Now()
	_, err := session.Execute(queryCtx, "SELECT * FROM customer", nil)
	end := time.Now()
	return end.Sub(start), err
}
