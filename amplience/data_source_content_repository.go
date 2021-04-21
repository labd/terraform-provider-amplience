package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContentRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Content Repositories function as subfolders inside of Hubs. Although a user can view content in " +
			"all repositories within a single hub, their ability to create content may be limited to certain " +
			"repositories. Typically you will want your content producers to be able to create content in one or more " +
			"repositories, but your planners to only be able to view the content. Content and slot types are registered " +
			"with hubs and enabled on repositories. So you can choose which types of content can be created in each " +
			"repository, or just choose to limit the number of content types that are available.\n" +
			"For more info see [Amplience Hubs & Repositories Docs](https://amplience.com/docs/intro/hubsandrepositories.html)",
		ReadContext: dataSourceContentRepositoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceContentRepositoryRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	repository_id := data.Get("id").(string)
	repository, err := ci.client.ContentRepositoryGet(repository_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(repository.ID)
	data.Set("label", repository.Label)
	data.Set("name", repository.Name)
	return diags
}
