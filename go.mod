module github.com/labd/terraform-provider-amplience

go 1.15

replace github.com/labd/amplience-go-sdk => /Users/mvantellingen/projects/labdigital/amplience-go-sdk

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.0
	github.com/labd/amplience-go-sdk v0.0.0-20210417163432-877c0ff03091
	github.com/stretchr/testify v1.6.1
)
