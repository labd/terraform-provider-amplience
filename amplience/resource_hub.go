package amplience

import (
	"context"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHub() *schema.Resource {
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
		CreateContext: resourceHubCreate,
		ReadContext:   resourceHubRead,
		UpdateContext: resourceHubUpdate,
		DeleteContext: resourceHubDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceHubCreate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	hub, err := ci.client.HubGet(ci.hubID)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceHubSaveState(data, hub)
	return diags
}

func resourceHubRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	hubId := data.Id()
	hub, err := ci.client.HubGet(hubId)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceHubSaveState(data, hub)
	return diags
}

func resourceHubUpdate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()

	instance, err := ci.client.HubGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	input := createResourceHubInput(data)
	hub, err := ci.client.HubUpdate(instance, input)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceHubSaveState(data, hub)
	return diags
}

func createResourceHubInput(data *schema.ResourceData) content.HubInput {
	input := content.HubInput{
		Label:       data.Get("label").(string),
		Description: data.Get("description").(string),
	}
	return input
}

// The amplience API does not have a hub delete functionality. Setting ID to "" and returning nil
func resourceHubDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceHubSaveState(data *schema.ResourceData, resource content.Hub) {
	data.SetId(resource.ID)
	_ = data.Set("label", resource.Label)
	_ = data.Set("description", resource.Description)

}
