local_test_version = 99.0.0

build-local:
	go build -o terraform-provider-amplience_${local_test_version}
	cp terraform-provider-amplience_${version} ~/.terraform.d/plugins/registry.terraform.io/hashicorp/amplience/${local_test_version}/darwin_amd64/terraform-provider-amplience_v${local_test_version}

format:
	go fmt ./...

test:
	go test -v ./...

deps:
	go mod tidy
	go mod vendor
