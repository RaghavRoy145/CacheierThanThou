# CachierThanThou

* Implements a Distributed Cache, multiple servers can join a "cluster" and clients can Set or Get from them by invoking the server nodes.
* Only the leader node is allowed to SET on member nodes.
* GET can be done on any node.

## Usage
* ```make run``` to start the leader server
* ```make runfollower``` to start follower nodes
* The main func runs a SET and GET from a client, this is just to see how it works
* To set listener address (```make runfollower``` implicitly sets this to :4000 and ```make run``` sets it to :3000) use the "listenaddr" flag
* For followers, to set leader address (```make runfollower``` implicitly sets this to :3000) use the "leaderaddr" flag

## TODO
* Currently, Raft consensus is only being implemented using mock servers (that can't dial each other)
* To run this: ```go run client/runtest/main.go```
* Need to apply this to the current server/client set up to make it real world
* Issue open, feel free to contribute
