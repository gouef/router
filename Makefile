.PHONY: install tests

install:
	go mod tidy && go mod vendor

tests:
	go test -v -covermode=set ./... -coverprofile=coverage.txt && go tool cover -func=coverage.txt