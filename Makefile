.PHONY: build test lint format

.DEFAULT:
build: test
	@goreleaser release --snapshot --clean --skip=publish

test:
	@go test -mod=vendor --cover -covermode=count -coverprofile=coverage.out ./...

lint:
	@golangci-lint run ./... -v

.ONESHELL:
format:
	@gofmt -l -e cmd/**/*.go
	@gofmt -l -e *.go
	@go fmt -mod=vendor ./...

update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif