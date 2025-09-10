binary_name = gameoflife

GOOS     := linux windows darwin
GOARCH   := amd64 arm64

PLATFORMS := $(foreach GOOS,$(GOOS),$(foreach GOARCH,$(GOARCH),$(GOOS)-$(GOARCH)))

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
	go build -o ./bin/${binary_name} main.go

.PHONY: run
run: build
	./bin/${binary_name}

