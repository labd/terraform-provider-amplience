package amplience_provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContentType_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContentTypeConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "content_type_uri", "http://www.example.com/content-type.json",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.#", "1",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.#", "2",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.0.size", "12",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.0.url", "http://www.example.com",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.1.size", "34",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.1.url", "http://www.google.com",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.#", "2",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.0.label", "label_1",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.0.templated_uri", "http://www.example.com",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.0.default", "false",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.1.label", "label_2",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.1.templated_uri", "http://www.google.com",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.1.default", "true",
					),
				),
			},
			{
				Config: testAccCheckContentTypeUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "content_type_uri", "http://www.example.com/content-type.json",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.#", "1",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.#", "1",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.0.size", "1234",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.icon.0.url", "http://www.example.com",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.#", "2",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.0.label", "label_1",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.0.templated_uri", "http://www.example.com/new-url",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.0.default", "true",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.1.label", "label_2",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.1.templated_uri", "http://www.example.com/new-url",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_type.test_type", "settings.0.visualization.1.default", "false",
					),
				),
			},
		},
	})
}

func testAccContentTypeConfig() string {
	return fmt.Sprint(`
	resource "amplience_content_type" "test_type" {
      content_type_uri = "https://schema-examples.com/blogpost"

      settings {
		label = "test_settings"
		
		icon {
          size = 12
          url = "http://www.example.com"
        } 

        icon {
          size = 34
          url = "http://www.google.com"
        }
		
        visualization {
          label = "label_1"
          templated_uri = "http://www.example.com"
          default = false
        }

        visualization {
          label = "label_2"
          templated_uri = "http://www.google.com"
          default = true
        }
      }
    }
`)
}

func testAccCheckContentTypeUpdate() string {
	return fmt.Sprint(`
	resource "amplience_content_type" "test_type" {
      content_type_uri = "https://schema-examples.com/blogpost"

      settings {
		label = "test_settings"
		
		icon {
          size = 1234
          url = "http://www.example.com"
        }
		
        visualization {
          label = "label_1"
          templated_uri = "http://www.example.com/new-url"
          default = true
        }

        visualization {
          label = "label_2"
          templated_uri = "http://www.google.com/new-url"
          default = false
        }
      }
    }
`)
}

func testAccCheckContentTypeDestroy(state *terraform.State) error {
	return nil
}
