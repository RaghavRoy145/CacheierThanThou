package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/raghavroy145/DistributedCaching/client"
)

type Server struct {
	raft *raft.Raft
}

func main() {
	var (
		cfg           = raft.DefaultConfig()
		fsm           = &raft.MockFSM{}
		logStore      = raft.NewInmemStore()
		stableStore   = raft.NewInmemStore()
		snapShotStore = raft.NewInmemSnapshotStore()
		timeout       = time.Second * 5
		// logOutout     = new(bytes.Buffer)
	)
	cfg.LocalID = "royboi"
	// copied from hashicorp raft test code
	ips, err := net.LookupIP("localhost")
	if err != nil {
		log.Fatal(err)
	}
	if len(ips) == 0 {
		log.Fatal("localhost did not resolve to any IPs")
	}
	addr := &net.TCPAddr{IP: ips[0], Port: 4000}
	tr, err := raft.NewTCPTransport(":4000", addr, 10, timeout, os.Stdout)
	if err != nil {
		log.Fatal("TCP net failed", err)
	}
	r, err := raft.NewRaft(cfg, fsm, stableStore, logStore, snapShotStore, tr)
	if err != nil {
		log.Fatal("Failed to create new raft: ", err)
	}
	fmt.Printf("%+v\n", r)
	select {}
}

func SendStuff() {
	c, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		var (
			key   = []byte(fmt.Sprintf("key_%d", i))
			value = []byte(fmt.Sprintf("val_%d", i))
		)

		err = c.Set(context.Background(), key, value, 0)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
	c.Close()
}
