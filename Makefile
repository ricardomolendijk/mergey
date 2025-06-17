# Makefile for mergeplease Go app

.PHONY: build test run clean

build:
	cd cmd && go build -o ../mergeplease main.go

test:
	go test ./...

run: build
	./mergeplease

clean:
	rm -f mergeplease
