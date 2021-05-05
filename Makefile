build:
	go build -o bin/sequence-length ./cmd/sequence-length/main.go
	go build -o bin/sequence-random ./cmd/sequence-random/main.go
	go build -o bin/sequence-shuffle ./cmd/sequence-shuffle/main.go
	go build -o bin/dataset-onehot ./cmd/dataset-onehot/main.go
	go build -o bin/kmer-count ./cmd/kmer-count/main.go