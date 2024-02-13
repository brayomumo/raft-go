kv-build:
	go build -o bin/kv src/store/*.go

kv-run: kv-build
	./bin/kv