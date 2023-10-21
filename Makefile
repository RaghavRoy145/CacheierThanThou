build:
	go build -o bin/DistributedCache

run: build
	./bin/DistributedCache