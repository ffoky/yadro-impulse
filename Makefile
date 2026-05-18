.PHONY: test
test:
	@go test  ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: lint
lint:
	@golangci-lint run -j8 --fix

.PHONY: format
format: fmt vet lint
