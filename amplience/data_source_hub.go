package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHub() *schema.Resource {
	return &schema.Resource{
		Description: "Permissions are set at the hub level. All users of a hub can at least view all of the " +
			"content within the repositories inside that hub. Content cannot be shared across hubs. However, content " +
			"can be shared and linked to across repositories within the same hub. So you can create a content item " +
			"in one repository and include content stored in another. Events and editions are scheduled within a " +
			"single hub. So if you want an overall view of the planning calendar across many brands, then you may wish " +
			"to consider a single hub. However, in some cases you may want to keep the calendars separate. Many " +
			"settings, such as the publishing endpoint (the subdomain to which your content is published) are set at " +
			"a hub level. Multiple hubs may publish content to the same endpoint.\n" +
			"For more info see [Amplience Hubs & Repositories Docs](https://amplience.com/docs/intro/hubsandrepositories.html)",
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
	ci := getClient(meta)

	hub_id := data.Get("id").(string)

	hub, err := ci.client.HubGet(hub_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(hub.ID)
	data.Set("hub_id", hub.ID)
	return diags
}
