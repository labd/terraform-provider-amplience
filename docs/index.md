# Amplience Terraform provider
This is the Terraform provider for the Amplience Dynamic Content API. It allows you to configure your
Amplience project with infrastructure-as-code principles.


# Commercial support
Need support implementing this terraform module in your organization? We are
able to offer support. Please contact us at
[opensource@labdigital.nl](opensource@labdigital.nl)!


## Installation
Terraform 0.13 added support for automatically downloading providers from
the terraform registry. Add the following to your terraform project

```hcl
terraform {
  required_providers {
    commercetools = {
      source = "labd/amplience"
    }
  }
}
```

Packages of the releases are available at [the GitHub Repo](https://github.com/labd/terraform-provider-amplience/releases).
See the [terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for more information about installing third-party providers.


## Using the provider
To use the provider you can either use environment variables, or directly configure the necessary parameters
on the provider itself

#### Environment Variables
The provider will read these values from the environment:
- `AMPLIENCE_CLIENT_ID`
- `AMPLIENCE_CLIENT_SECRET`
- `AMPLIENCE_HUB_ID`
- `AMPLIENCE_CONTENT_API_PATH`

#### Terraform Config
In order to configure the provider directly set the following fields
```hcl
provider "amplience" {
  client_id        = "<your client id>"
  client_secret    = "<your client secret>"
  hub_id           = "<your hub id"
  content_api_path = "<your content api path>"
}
```

## Authors
This project is developed by [Lab Digital](https://www.labdigital.nl). We
welcome additional contributors. Please see our
[GitHub repository](https://github.com/labd/terraform-provider-commercetools)
for more information.
