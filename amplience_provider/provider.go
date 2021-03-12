package amplience_provider

import (
	"context"
	"fmt"

	"github.com/labd/terraform-provider-amplience/amplience"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CLIENT_ID", nil),
				Description: "The OAuth Client ID for the Amplience management API https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication",
				Sensitive:   true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CLIENT_SECRET", nil),
				Description: "The OAuth Client Secret for Amplience management API. https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication",
				Sensitive:   true,
			},
			"hub_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_HUB_ID", nil),
				Description: "The Hub ID of the Amplience Hub to use this provider instance with",
				Sensitive:   false,
			},
			"content_api_path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CONTENT_API_PATH", "https://api.amplience.net/v2/content"),
				Description: "The base URL path for the Amplience Content API",
				Sensitive:   false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"amplience_content_repository": resourceContentRepository(),
			"amplience_webhook":            resourceWebhook(),
		},
		ConfigureContextFunc: amplienceProviderConfigure,
	}
}

// amplienceProviderConfigure should instantiate an Amplience client from the env vars when a proper Amplience client
// library is implemented. For now it just sets the values to a "client" struct for further use
func amplienceProviderConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	hubID := data.Get("hub_id").(string)
	if hubID == "" {
		return nil, diag.FromErr(fmt.Errorf("hub_id is empty, can not instantiate provider"))
	}

	clientID := data.Get("client_id").(string)
	if clientID == "" {
		return nil, diag.FromErr(fmt.Errorf("client_id is empty, can not instantiate provider"))
	}
	clientSecret := data.Get("client_secret").(string)
	if clientSecret == "" {
		return nil, diag.FromErr(fmt.Errorf("client_secret is empty, can not instantiate provider"))
	}
	contentAPIPath := data.Get("content_api_path").(string)

	client := &amplience.Client{
		ID:             clientID,
		Secret:         clientSecret,
		HubID:          hubID,
		ContentAPIPath: contentAPIPath,
	}
	return client, diags
}
