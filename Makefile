INSTALL_DIR = /usr/local/bin

ifdef prefix
	INSTALL_DIR = $(prefix)
endif

build:
	go build -o bin/dataset-check-duplicated ./cmd/dataset-check-duplicated/main.go
	go build -o bin/dataset-onehot ./cmd/dataset-onehot/main.go
	go build -o bin/dataset-split ./cmd/dataset-split/main.go
	go build -o bin/fastq-filter ./cmd/fastq-filter/main.go

install:
	cp bin/dataset-check-duplicated $(INSTALL_DIR)/dataset-check-duplicated
	cp bin/dataset-onehot $(INSTALL_DIR)/dataset-onehot
	cp bin/dataset-split $(INSTALL_DIR)/dataset-split
	cp bin/fastq-filter $(INSTALL_DIR)/fastq-filter

uninstall:
	rm -f $(INSTALL_DIR)/dataset-check-duplicated
	rm -f $(INSTALL_DIR)/dataset-onehot
	rm -f $(INSTALL_DIR)/dataset-split
	rm -f $(INSTALL_DIR)/fastq-filter