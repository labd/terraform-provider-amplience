package amplience_provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccWebhooks_createAndUpdate(t *testing.T) {
	webhookLabel := acctest.RandomWithPrefix("webhook-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccWebhooksDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhooksConfig(webhookLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "label", webhookLabel),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.#", "2"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.0", "dynamic-content.content-item.created"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.1", "dynamic-content.content-item.updated"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "handlers.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "active", "false"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "secret", "a-test-secret"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "notifications.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "notifications.0.email", "example.person@gmail.com"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.#", "2"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.0.key", "X-Additional-Header"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.0.value", "abc123"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.1.key", "X-second-Header"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.1.secret_value", "321cba"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.#", "2"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.type", "equal"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.0.json_path", "$.payload.id"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.0.value.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.0.value.0", "abc"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.type", "in"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.json_path", "$.payload.id"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.#", "2"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.0", "abc"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.1", "123"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "method", "POST"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "custom_payload.type", "text/x-handlebars-template"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "custom_payload.value", "OPEN_INVERSE"),
				),
			},
			{
				Config: testAccWebhookUpdate(webhookLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "label", webhookLabel),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.#", "3"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.0", "dynamic-content.content-item.created"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.1", "dynamic-content.content-item.updated"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "events.2", "dynamic-content.content-item.workflow.updated"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "handlers.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "active", "true"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "secret", "an-updated-test-secret"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "notifications.#", "1"),

					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "notifications.0.email", "example.person@gmail.com"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.#", "2"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.0.key", "X-Additional-Header"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.0.value", "abc123"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.1.key", "X-second-Header"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "header.1.secret_value", "321cba"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.#", "2"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.type", "equal"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.0.json_path", "$.payload.id"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.0.value.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.0.arguments.0.value.0", "123"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.type", "in"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.#", "1"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.json_path", "$.payload.id"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.#", "3"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.0", "abc"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.1", "123"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "filter.1.arguments.0.value.2", "a third updated value"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "method", "PATCH"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "custom_payload.type", "text/x-handlebars-template"),
					resource.TestCheckResourceAttr(
						"amplience_webhook.standard", "custom_payload.value", "OPEN_INVERSE"),
				),
			},
		},
	})
}

func testAccWebhooksConfig(label string) string {
	return fmt.Sprintf(`
	resource "amplience_webhook" "standard" {
      label = "%[1]s"

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
      }
	  
      header {
       key = "X-second-Header"
       secret_value = "321cba"
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


    }`, label)
}

func testAccWebhookUpdate(label string) string {
	return fmt.Sprintf(`
    resource "amplience_webhook" "standard" {
      label = "%[1]s"

      events = [
		"dynamic-content.content-item.created",
		"dynamic-content.content-item.updated",
		"dynamic-content.content-item.workflow.updated",
      ]
      handlers = [
		"http://example.com/webhook",
      ]
      active   = true

      notifications {
		email = "example.person@gmail.com"
	  }

      secret = "an-updated-test-secret"

	  header {
       key = "X-Additional-Header"
       value = "abc123"
      }
	  
      header {
       key = "X-second-Header"
       secret_value = "321cba"
      }

      filter {
       type = "equal"
       arguments {
         json_path = "$.payload.id"
         value = ["123"]
       }
      }
      
      filter {
       type = "in"
       arguments {
         json_path = "$.payload.id"
         value = ["abc", "123", "a third updated value"]
       }
      }

      method = "PATCH"

      custom_payload = {
        type = "text/x-handlebars-template"
        value = "OPEN_INVERSE"
      }

    }`, label)

}

// TODO: Implement
func testAccWebhooksDestroy(s *terraform.State) error {
	return nil
}
