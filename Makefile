build:
	CGO_ENABLED=0 go build -o flume ./cmd/flume/

run: build
	./flume

test:
	go test ./...

clean:
	rm -f flume

.PHONY: build run test clean
