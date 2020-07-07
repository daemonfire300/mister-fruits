GOPATH?=~/go

lint: $(GOPATH)/bin/golangci-lint
	$(GOPATH)/bin/golangci-lint run

$(GOPATH)/bin/golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.27.0

