run: build 
	./omnisync

build: 
	go build ./cmd/omnisync

test:
	go test ./... -v

clean: 
	rm -rf ./cmd/omnisync/main
	rm -rf ./omnisync	