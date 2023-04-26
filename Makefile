.DEFAULT_GOAL := test

.PHONY: clean test vet

vet:
	go vet ./...

test:
	go test -v ./...

clean:
	go clean


