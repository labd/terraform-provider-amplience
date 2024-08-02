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
