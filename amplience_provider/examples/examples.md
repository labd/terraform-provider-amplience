# Resources
## Webhook resource
```terraform
resource "amplience_webhook" "standard" {
  label = "test_webhook_label"

  events   = ["dynamic-content.content-item.created", "dynamic-content.content-item.updated"]
  handlers = ["http://test-url.com/webhook"]
  active   = true

  notifications = {
    email = "test@example.com"
  }

  secret = "a-test-secret"
  headers = {
    key    = "X-additional-Header"
    value  = "testval123"
    secret = true
  }

// TODO:
//  filters =
//

  method = "POST"
  custom_payload = {
    type  = "text/x-handlebars-template"
    value = "{{#withDeliveryContentItem contentItemId=payload.rootContentItem.id account=\"account-name\" stagingEnvironment=\"staging-environment-url\"}} ... {/withDeliveryContentItem}}"
  }
}
```
