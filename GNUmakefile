GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: fmtcheck
	go install

init:
	go build -o terraform-provider-vcda
	terraform init

plan: init
	terraform plan

apply: init
	terraform apply

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

testacc:
	TF_ACC=1 go test -v -run $(TESTS) -timeout 10m ./...

.PHONY: build init plan apply fmt fmtcheck testacc