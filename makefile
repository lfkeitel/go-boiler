.PHONY: build doc fmt install lint test vet

default: build

build: vet
	go build -v -o ./bin/boiler ./src/cmd/boiler

doc:
	godoc -http=:6060 -index

fmt:
	go fmt ./src/...

install: vet
	go install -v ./src/cmd/./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

test:
	go test ./src/...

vet:
	go vet ./src/...
