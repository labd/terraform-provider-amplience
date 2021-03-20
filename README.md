# Amplience Terraform Provider
Terraform provider for [Amplience](https://amplience.com/).

The intention of this provider is to cover the [Amplience dynamic content management APIs](https://amplience.com/docs/api/dynamic-content/management/index.html), so that one can manage an entire Amplience configuration through Terraform.

One provider can manage the resource of one HubID

## Currently supported resources

Currently the checked resources are supported. Support for additional resources will come when they are required in projects, or contributed.

- [ ] Administration
- [ ] Content Items
- [x] Content Repositories
- [ ] Content Types
- [ ] Editions
- [ ] Events
- [ ] Extensions
- [ ] Folders
- [ ] Hubs
- [ ] Localization
- [ ] Integrations
- [ ] Publishing Jobs
- [ ] SFCC
- [ ] SFMC
- [ ] Search Indexes
- [ ] Search Indexes - Analyics
- [ ] Slots
- [ ] Snapshots
- [x] Webhooks
- [ ] Workflows
- [ ] Hierarchy Node

# Installation

## Terraform registry

Terraform 0.13 added support for automatically downloading providers from
the terraform registry. Simply add the following to your terraform project to use the latest release

```hcl
terraform {
  required_providers {
    amplience = {
      source = "labd/amplience"
    }
  }
}
```

## Binaries

Packages of the releases are available at
https://github.com/labd/terraform-provider-amplience/releases See the
[terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for more information about installing third-party providers.

# Getting started

[Read our documentation](https://registry.terraform.io/providers/labd/amplience/latest/docs).

# Contributing

## Building the provider
Clone repository to: `$GOPATH/src/github.com/labd/terraform-provider-amplience` 

Then run
```sh
make build
```

### Testing local changes
As of terraform 0.13 testing local changes requires a little effort.
You can run 
```sh
make build-local
```
To build the provider with a very high version number and copy it to your terraform plugins folder (default is for Mac, 
change OS_ARCH if running Linux or change path if running Windows)
If you set your provider source as `labd/amplience` and the version to your built `version` it should use the local
provider. See also [the Terraform 0.13 upgrade guide](https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers)


## Debugging / Troubleshooting

There is currently one environment settings for troubleshooting:

- `TF_LOG=1` enables debug output for Terraform.

Note this generates a lot of output!

## Testing

### Running the unit tests

```sh
$ make test
```

### Running an Acceptance Test

In order to run the full suite of Acceptance tests, run `make testacc`.

**NOTE:** Acceptance tests create real resources.

Prior to running the tests provider configuration details such as access keys
must be made available as environment variables.

Since we need to be able to create Amplience resources, we need the
Amplience API credentials. So in order for the acceptance tests to run
correctly please provide all of the following:

```sh
export AMPLIENCE_CLIENT_ID=...
export AMPLIENCE_CLIENT_SECRET=...
export AMPLIENCE_HUB_ID=...
```

For convenience, place a `testenv.sh` in your `local` folder (which is
included in .gitignore) where you can store these environment variables.

Tests can then be started by running

```sh
$ source local/testenv.sh
$ make testacc
```

## Releasing

When pushing a new tag prefixed with `v` a GitHub action will automatically
use Goreleaser to build and release the build.

```sh
git tag <release> -m "Release <release>" # please use semantic version, so always vX.Y.Z
git push --follow-tags
```

## TODO List 
- Currently this repository contains a minimal `amplience` package to call the Amplience API with. However one of the first things 
that should be improved for this project is that a proper Amplience Client Library should be set up
- The above Client Library should implement support for multiple Hubs in a clear and logical manner (1 Hub = 1 Provider)
- Unit/acceptance tests should be expanded
- It would be nice to have a Mock Amplience server to run (non-acceptance) tests against
- The above tests can then be made to run on push through a Github Action

## Authors

This project is developed by [Lab Digital](https://www.labdigital.nl). We
welcome additional contributors. Please see our
[GitHub repository](https://github.com/labd/terraform-provider-amplience)
for more information.
