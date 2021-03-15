package amplience_provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContentRepository_CreateAndUpdate(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test-repo")
	label := acctest.RandomWithPrefix("tf-acc-test-repo")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContentRepositoryConfig(name, label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"amplience_content_repository.testrepo", "name", name,
					),
					resource.TestCheckResourceAttr(
						"amplience_content_repository.testrepo", "label", label,
					),
				),
			},
			{
				Config: testAccContentRepositoryConfig(name, label+"-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"amplience_content_repository.testrepo", "name", name+"-updated",
					),
					resource.TestCheckResourceAttr(
						"amplience_content_repository.testrepo", "label", label+"-updated",
					),
				),
			},
		},
	})
}

func testAccContentRepositoryConfig(name, label string) string {
	return fmt.Sprintf(`
	resource "amplience_content_repository" "testrepo" {
      name = "%[1]s"
      label = "%[2]s"
    }
`, name, label)
}

func testAccCheckContentRepositoryDestroy(state *terraform.State) error {
	//TODO: Implement
	return nil
}
