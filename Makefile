build:
	go build -o bin/DistributedCache

run: build
	./bin/DistributedCache

runfollower: build
	./bin/DistributedCache --listenaddr :4000 --leaderaddr :3000