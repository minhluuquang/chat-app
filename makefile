.PHONY: build dev

build:
	go build -o server chat/*.go

dev: build
	./server; rm server