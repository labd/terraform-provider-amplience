.PHONY: build, build-local, format, test, testacc, deps, docs
LOCAL_TEST_VERSION = 99.0.0
OS_ARCH = darwin_amd64

build:
	go build

# Build local provider with very high version number for easier local testing and debugging
# see: https://discuss.hashicorp.com/t/easiest-way-to-use-a-local-custom-provider-with-terraform-0-13/12691/5
build-local:
	go build -o terraform-provider-amplience_${LOCAL_TEST_VERSION}
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/labd/amplience/${LOCAL_TEST_VERSION}/${OS_ARCH}
	cp terraform-provider-amplience_${LOCAL_TEST_VERSION} ~/.terraform.d/plugins/registry.terraform.io/labd/amplience/${LOCAL_TEST_VERSION}/${OS_ARCH}/terraform-provider-amplience_v${LOCAL_TEST_VERSION}

format:
	go fmt ./...

# Added grep is to remove clutter in the test logs from the parallelisation
test:
	go test -v ./... | grep -v 'PAUSE\|CONT'

testacc:
	TF_ACC=1 go test -v ./...

deps:
	go mod tidy
	go mod vendor

docs:
	tfplugindocs

update-sdk:
	GO111MODULE=on go get github.com/labd/amplience-go-sdk
	GO111MODULE=on go mod vendor
	GO111MODULE=on go mod tidy
