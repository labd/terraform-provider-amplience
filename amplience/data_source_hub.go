package amplience

import (
	"context"

	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHub() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHubRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceHubRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id := data.Get("id").(string)

	hub, err := c.HubGet(hub_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(hub.ID)
	data.Set("hub_id", hub.ID)
	return diags
}
