.PHONY: pre-commit
pre-commit:
	go mod tidy
	gofumpt -l -w .
	golangci-lint run .

.PHONY: test
test:
	go test .
	go test -v -race -count=100 .

.PHONY: clean
	go clean -cache -modcache -i -r

.DEFAULT_GOAL := pre-commit