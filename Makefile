binary_name = gameoflife

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

.PHONY: clean
clean:
	go clean
	/bin/rm -rf ./bin/

.PHONY: build
build:
	go build -o=./bin/${binary_name} .

.PHONY: run
run: build
	./bin/${binary_name}

