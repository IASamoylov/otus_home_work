.PHONY: pre-commit
pre-commit:
	go mod tidy
	gofumpt -l -w .
	gci -w .
	golangci-lint run .

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