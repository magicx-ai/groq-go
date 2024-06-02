lint:
	golangci-lint run

format:
	go mod tidy
	golangci-lint run --fix

test:
	go test -cover -v ./...

integration-test:
	go test -tags=integration -v ./...

generate:
	@go install go.uber.org/mock/mockgen@latest
	go generate ./...

