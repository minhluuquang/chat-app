.PHONY: build dev

build:
	go build -o server chat/*.go

dev:
	./server