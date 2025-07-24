GOBIN ?= $$(go env GOPATH)/bin

.PHONY: install-go-test-coverage
install-go-test-coverage:
	@if [ ! -f "${GOBIN}/go-test-coverage" ]; then \
		echo "Installing go-test-coverage..."; \
		go install github.com/vladopajic/go-test-coverage/v2@latest; \
	fi

.PHONY: test
test:
	@echo "Running all tests..."
	go test ./...

.PHONY: coverage
coverage: install-go-test-coverage
	go test ./... -coverprofile=./coverage.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yaml --profile=coverage.out

.PHONY: lint
lint:
	 golangci-lint run
