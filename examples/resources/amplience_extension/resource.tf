resource "amplience_extension" "example" {
  name                          = "extension-01"
  label                         = "Extension-01 first-draft"
  description                   = "This is a description of an extension."
  url                           = "https://example.com"
  height                        = 400
  enabled_for_all_content_types = false
  category                      = "CONTENT_FIELD"

  settings = jsonencode({
    API = {
      READ = true
      EDIT = true
    }
    SANDBOX = {
      SAME_ORIGIN          = true
      MODALS               = false
      NAVIGATION           = false
      POPUPS               = true
      POPUP_ESCAPE_SANDBOX = false
      DOWNLOADS            = false
      FORMS                = false
    }
  })

  parameters = jsonencode({
    color = "blue"
  })

  snippets = [
    {
      label = "my snippet"
      body  = jsonencode({})
    }
  ]
}
