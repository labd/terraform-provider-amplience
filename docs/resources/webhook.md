---

# amplience_webhook (Resource)

A webhook is a way for Dynamic Content to automatically send messages or data to a third party system.
For example, an e-commerce system might need to know when an edition is scheduled to allow it to retrieve the slots and
content that the edition contains. Similarly it would also need to be notified when that edition is unscheduled.

This provider allows full CRUD functionality for Amplience webhooks

For more information see the Amplience [user documentation](https://amplience.com/docs/integration/webhooks.html)
and/or the [API documentation](https://amplience.com/docs/api/dynamic-content/management/index.html#tag/Webhooks)

## Important
The response from the API when creating or updating webhooks contains secrets (the `email` field in `notifications` and 
the `value` field in `header`) which are returned as `null`.

As this can lead to state issues for Terraform the provider has been configured so that a **new** Webhook is created
when these fields are changed. If only other fields are changed the existing Webhook will be altered.

## Example Usage
```hcl
	resource "amplience_webhook" "standard" {
  label = "webhook_example_label"

  events = [
    "dynamic-content.content-item.created",
    "dynamic-content.content-item.updated",
  ]
  handlers = [
    "http://example.com/webhook",
  ]

  notifications {
    email = "example.person@gmail.com"
  }

  header {
    key = "X-Additional-Header"
    value = "abc123"
    secret = false
  }

  header {
    key = "X-second-Header"
    value = "321cba"
    secret = true
  }

  filter {
    type = "equal"
    arguments {
      json_path = "$.payload.id"
      value = ["abc"]
    }
  }

  filter {
    type = "in"
    arguments {
      json_path = "$.payload.id"
      value = ["abc", "123"]
    }
  }

  active   = false
  secret = "a-test-secret"

  method = "POST"

  custom_payload = {
    type = "text/x-handlebars-template"
    value = "OPEN_INVERSE"
  }
}
```