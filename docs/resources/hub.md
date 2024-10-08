---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "amplience_hub Resource - terraform-provider-amplience"
subcategory: ""
description: |-
  Permissions are set at the hub level. All users of a hub can at least view all of the content within the repositories inside that hub. Content cannot be shared across hubs. However, content can be shared and linked to across repositories within the same hub. So you can create a content item in one repository and include content stored in another. Events and editions are scheduled within a single hub. So if you want an overall view of the planning calendar across many brands, then you may wish to consider a single hub. However, in some cases you may want to keep the calendars separate. Many settings, such as the publishing endpoint (the subdomain to which your content is published) are set at a hub level. Multiple hubs may publish content to the same endpoint.
  For more info see Amplience Hubs & Repositories Docs https://amplience.com/docs/intro/hubsandrepositories.html
  It is recommended to import a new hub instead of creating it! This is because the hub already exists, so any differences in configuration might be overwritten, leading to unintended outcomes.
---

# amplience_hub (Resource)

Permissions are set at the hub level. All users of a hub can at least view all of the content within the repositories inside that hub. Content cannot be shared across hubs. However, content can be shared and linked to across repositories within the same hub. So you can create a content item in one repository and include content stored in another. Events and editions are scheduled within a single hub. So if you want an overall view of the planning calendar across many brands, then you may wish to consider a single hub. However, in some cases you may want to keep the calendars separate. Many settings, such as the publishing endpoint (the subdomain to which your content is published) are set at a hub level. Multiple hubs may publish content to the same endpoint.

For more info see [Amplience Hubs & Repositories Docs](https://amplience.com/docs/intro/hubsandrepositories.html)

**It is recommended to import a new hub instead of creating it!** This is because the hub already exists, so any differences in configuration might be overwritten, leading to unintended outcomes.

## Example Usage

```terraform
resource "amplience_hub" "my-hub" {
  name  = "myhub"
  label = "My Hub"
  settings = {
    applications = [
      {
        name         = "Application"
        template_uri = "https://application.com/preview/"
      }
    ]
    asset_management = {
      client_config = "HUB"
      enabled       = true
    }
    devices = [
      {
        name      = "Desktop"
        width     = 1024
        height    = 1024
        orientate = false
      },
      {
        name      = "Tablet"
        width     = 640
        height    = 768
        orientate = false
      },
      {
        name      = "Mobile"
        width     = 320
        height    = 512
        orientate = false
      }
    ]
    localization = {
      locales = [
        "en-GB",
        "nl-NL"
      ]
    }
    preview_virtual_staging_environment = {
      hostname = "my-preview-virtual-staging.environment.io"
    }
    virtual_staging_environment = {
      hostname = "my-virtual-staging.environment.io"
    }
    publishing = {
      platforms = {
        amplience_dam = {
          api_key    = "my-api-key"
          api_secret = "my-api-secret"
          endpoint   = "my-endpoint"
        }
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `label` (String) Hub label
- `name` (String) Hub name

### Optional

- `description` (String) Hub description
- `settings` (Attributes) Hub settings (see [below for nested schema](#nestedatt--settings))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `applications` (Attributes List) (see [below for nested schema](#nestedatt--settings--applications))
- `asset_management` (Attributes) (see [below for nested schema](#nestedatt--settings--asset_management))
- `devices` (Attributes List) (see [below for nested schema](#nestedatt--settings--devices))
- `localization` (Attributes) (see [below for nested schema](#nestedatt--settings--localization))
- `preview_virtual_staging_environment` (Attributes) (see [below for nested schema](#nestedatt--settings--preview_virtual_staging_environment))
- `publishing` (Attributes) (see [below for nested schema](#nestedatt--settings--publishing))
- `virtual_staging_environment` (Attributes) (see [below for nested schema](#nestedatt--settings--virtual_staging_environment))

<a id="nestedatt--settings--applications"></a>
### Nested Schema for `settings.applications`

Required:

- `name` (String)
- `template_uri` (String)


<a id="nestedatt--settings--asset_management"></a>
### Nested Schema for `settings.asset_management`

Optional:

- `client_config` (String)
- `enabled` (Boolean)


<a id="nestedatt--settings--devices"></a>
### Nested Schema for `settings.devices`

Required:

- `height` (Number)
- `name` (String)
- `orientate` (Boolean)
- `width` (Number)


<a id="nestedatt--settings--localization"></a>
### Nested Schema for `settings.localization`

Optional:

- `locales` (List of String)


<a id="nestedatt--settings--preview_virtual_staging_environment"></a>
### Nested Schema for `settings.preview_virtual_staging_environment`

Required:

- `hostname` (String) Virtual Staging Environment hostname


<a id="nestedatt--settings--publishing"></a>
### Nested Schema for `settings.publishing`

Optional:

- `platforms` (Attributes) (see [below for nested schema](#nestedatt--settings--publishing--platforms))

<a id="nestedatt--settings--publishing--platforms"></a>
### Nested Schema for `settings.publishing.platforms`

Optional:

- `amplience_dam` (Attributes) (see [below for nested schema](#nestedatt--settings--publishing--platforms--amplience_dam))

<a id="nestedatt--settings--publishing--platforms--amplience_dam"></a>
### Nested Schema for `settings.publishing.platforms.amplience_dam`

Required:

- `api_key` (String) DAM publishing client key
- `api_secret` (String, Sensitive) DAM publishing client secret
- `endpoint` (String) Publishing endpoint, also known as Company Tag




<a id="nestedatt--settings--virtual_staging_environment"></a>
### Nested Schema for `settings.virtual_staging_environment`

Required:

- `hostname` (String) Virtual Staging Environment hostname
