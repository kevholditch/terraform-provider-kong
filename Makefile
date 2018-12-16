GOFMT_FILES?=$$(find ./ -name '*.go' | grep -v vendor)

default: build test testacc

travisbuild: deps default

test:
	go test -v . ./kong

testacc:
	TF_ACC=1 go test -v ./kong -run="TestAcc"

build: goimportscheck vet testacc
	@go install
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-kong ~/.terraform.d/plugins/terraform-provider-kong
	@echo "Build succeeded"

build-gox: deps vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-kong" .

release:
	@curl -sL http://git.io/goreleaser | bash

deps:
	go get -u golang.org/x/net/context; \
    go get -u github.com/mitchellh/gox; \

clean:
	rm -rf pkg/

goimports:
	goimports -w $(GOFMT_FILES)

install-goimports:
	@go get golang.org/x/tools/cmd/goimports

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test testacc vet goimports goimportscheck errcheck vendor-status test-compile
