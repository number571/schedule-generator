CC=go build
.PHONY: default build clean
default: build
build: main.go
	$(CC) main.go
clean:
	rm -rf main
