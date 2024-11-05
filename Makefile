run: build 
	./omnisync

build: 
	go build ./cmd/omnisync

test:
	go test ./...

clean: 
	rm -rf ./cmd/omnisync/main
	rm -rf ./omnisync	