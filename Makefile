install:
	go install -ldflags="-X main.version=$(shell git describe --always --abbrev=6 --dirty=-dev)"
