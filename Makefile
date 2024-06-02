lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0
	golangci-lint run

format:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0
	go mod tidy
	golangci-lint run --fix

test:
	go test -cover -v ./...

integration-test:
	go test -tags=integration -v ./...

generate:
	@go install go.uber.org/mock/mockgen@latest
	go generate ./...

