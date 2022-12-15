build:
	go build -o bin/dataset-onehot ./cmd/dataset-onehot/main.go
	go build -o bin/dataset-split ./cmd/dataset-split/main.go
	go build -o bin/fastq-filter ./cmd/fastq-filter/main.go