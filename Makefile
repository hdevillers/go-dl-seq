build:
	go build -o bin/sequence-length ./cmd/sequence-length/main.go
	go build -o bin/sequence-random ./cmd/sequence-random/main.go
	go build -o bin/sequence-shuffle ./cmd/sequence-shuffle/main.go