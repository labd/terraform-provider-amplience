package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContentRepository() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentRepositoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
		},
	}
}

func dataSourceContentRepositoryRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	repository_id := data.Get("id").(string)
	repository, err := c.ContentRepositoryGet(repository_id)
	if err != nil {
		return diag.FromErr(err)
	}

	hub, err := repository.GetHub(c)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(repository.ID)
	data.Set("name", repository.Name)
	data.Set("label", repository.Label)
	data.Set("hub_id", hub.ID)
	return diags
}
