package main

import (
	"context"
	"fmt"
	"log"
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

	// mostly black box for me here lol, no clue why it works -> need to apply this to my server
	tr, err := raft.NewTCPTransport("127.0.0.1:4000", nil, 10, timeout, os.Stdout)
	if err != nil {
		log.Fatal("TCP net failed: ", err)
	}

	// could be a loop, sherlock
	server := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(cfg.LocalID),
		Address:  raft.ServerAddress("127.0.0.1:4000"),
	}

	server2 := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID("FOOBAR"),
		Address:  raft.ServerAddress("127.0.0.1:4001"),
	}

	server3 := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID("BARFOO"),
		Address:  raft.ServerAddress("127.0.0.1:4002"),
	}
	serverConfig := raft.Configuration{
		Servers: []raft.Server{server, server2, server3},
	}

	r, err := raft.NewRaft(cfg, fsm, stableStore, logStore, snapShotStore, tr)
	if err != nil {
		log.Fatal("Failed to create new raft: ", err)
	}

	r.BootstrapCluster(serverConfig)
	// raft.BootstrapCluster(cfg, logStore, logStore, snapShotStore, tr, serverConfig)

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
