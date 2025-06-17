# Makefile for mergeplease Go app

.PHONY: build test run clean

build:
	@if [ ! -f ~/.merge/config.yaml ]; then \
	  mkdir -p ~/.merge && cp config.yaml.example ~/.merge/config.yaml; \
	  echo "Created config.yaml at ~/.merge/config.yaml"; \
	fi
	cd cmd && go build -o ../mergeplease main.go


test:
	go test ./...

run: build
	@if [ ! -f ~/.merge/config.yaml ]; then \
	  mkdir -p ~/.merge && cp config.yaml.example ~/.merge/config.yaml; \
	fi
	./mergeplease

clean:
	rm -f mergeplease
