package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Content Repositories function as subfolders inside of Hubs. Although a user can view content in " +
			"all repositories within a single hub, their ability to create content may be limited to certain " +
			"repositories. Typically you will want your content producers to be able to create content in one or more " +
			"repositories, but your planners to only be able to view the content. Content and slot types are registered " +
			"with hubs and enabled on repositories. So you can choose which types of content can be created in each " +
			"repository, or just choose to limit the number of content types that are available.\n" +
			"For more info see [Amplience Hubs & Repositories Docs](https://amplience.com/docs/intro/hubsandrepositories.html)",
		CreateContext: resourceContentRepositoryCreate,
		ReadContext:   resourceContentRepositoryRead,
		UpdateContext: resourceContentRepositoryUpdate,
		DeleteContext: resourceContentRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceContentRepositoryCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	input := content.ContentRepositoryInput{
		Name:  data.Get("name").(string),
		Label: data.Get("label").(string),
	}

	repository, err := ci.client.ContentRepositoryCreate(ci.hubID, input)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(repository.ID)
	data.Set("name", repository.Name)
	data.Set("label", repository.Label)
	return diags
}

func resourceContentRepositoryRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	repository_id := data.Id()

	repository, err := ci.client.ContentRepositoryGet(repository_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("name", repository.Name)
	data.Set("label", repository.Label)
	return diags
}

func resourceContentRepositoryUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	repository_id := data.Id()

	if data.HasChange("label") || data.HasChange("name") {
		current, err := ci.client.ContentRepositoryGet(repository_id)
		if err != nil {
			return diag.FromErr(err)
		}

		input := content.ContentRepositoryInput{
			Name:  data.Get("name").(string),
			Label: data.Get("label").(string),
		}

		repository, err := ci.client.ContentRepositoryUpdate(current, input)
		if err != nil {
			return diag.FromErr(err)
		}

		data.Set("label", repository.Label)
		data.Set("name", repository.Name)
	}

	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentRepositoryDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	data.SetId("")
	return diags
}
