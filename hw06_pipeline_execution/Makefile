.PHONY: pre-commit
pre-commit:
	go mod tidy
	gofumpt -l -w .
	golangci-lint run .
	gci -w .

.PHONY: test-race
test-race:
	go test -v -race -count=100 .

.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean -cache -modcache -i -r

.DEFAULT_GOAL := pre-commit