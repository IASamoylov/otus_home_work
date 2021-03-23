.PHONY: pre-commit
pre-commit:
	gofumpt -l -w .
	golangci-lint run .

.DEFAULT_GOAL := pre-commit