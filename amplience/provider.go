package amplience

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/labd/amplience-go-sdk/content"
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
			"content_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CONTENT_API_URL", "https://api.amplience.net/v2/content"),
				Description: "The base URL path for the Amplience Content API",
				Sensitive:   false,
			},
			"auth_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_AUTH_URL", "https://auth.adis.ws/oauth/token"),
				Description: "The Amplience authentication URL",
				Sensitive:   false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"amplience_content_repository":      resourceContentRepository(),
			"amplience_content_type":            resourceContentType(),
			"amplience_content_type_assignment": resourceContentTypeAssignment(),
			"amplience_content_type_schema":     resourceContentTypeSchema(),
			"amplience_webhook":                 resourceWebhook(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"amplience_hub":                dataSourceHub(),
			"amplience_content_repository": dataSourceContentRepository(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	apiURL := d.Get("content_api_url").(string)
	authURL := d.Get("auth_url").(string)

	var diags diag.Diagnostics

	// FIXME: pass context to amplience sdk client
	spew.Dump(clientID)
	spew.Dump(clientSecret)

	client, err := content.NewClient(&content.ClientConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		URL:          apiURL,
		AuthURL:      authURL,
	})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}
