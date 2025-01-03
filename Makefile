.PHONY: install tests coverage

install:
	go mod tidy && go mod vendor

tests:
	go test -v -covermode=set ./... -coverprofile=coverage.txt && go tool cover -func=coverage.txt
coverage:
	go test -v -covermode=set ./... -coverprofile=coverage.txt && go tool cover -html=coverage.txt -o coverage.html && xdg-open coverage.html