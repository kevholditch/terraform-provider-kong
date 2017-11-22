GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build test testacc

test: fmtcheck
	go test -v . ./kong

testacc: fmtcheck
	go test -v ./kong -run="TestAcc"

build: fmtcheck vet testacc
	@go install
	@go get golang.org/x/net/context
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-kong ~/.terraform.d/plugins/terraform-provider-kong
	@echo "Build succeeded"

build-gox: deps fmtcheck vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-kong" .

release:
	go get github.com/goreleaser/goreleaser; \
    goreleaser; \

deps:
	go get -u github.com/mitchellh/gox

clean:
	rm -rf pkg/
fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile