.PHONY: pre-commit
pre-commit:
	go mod tidy
	gofumpt -l -w .
	golangci-lint run .

.PHONY: test
test:
	go test .

.DEFAULT_GOAL := pre-commit