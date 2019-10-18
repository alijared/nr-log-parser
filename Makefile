.PHONY: build
build: ## installs the application
	@go build -o ${GOPATH}/bin/nrlp
