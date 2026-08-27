.PHONY: build fmt test testacc

build:
	go build -o terraform-provider-ag5 .

fmt:
	gofmt -w $$(find . -name '*.go')
	terraform fmt -recursive examples

test:
	go test ./...

testacc:
	TF_ACC=1 go test -count=1 ./...
