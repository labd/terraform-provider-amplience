package amplience

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Description: "The OAuth Client ID for the Amplience management API https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CLIENT_ID", nil),
				Sensitive:   true,
			},
			"client_secret": {
				Description: "The OAuth Client Secret for Amplience management API. https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CLIENT_SECRET", nil),
				Sensitive:   true,
			},
			"content_api_url": {
				Description: "The base URL path for the Amplience Content API",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_CONTENT_API_URL", "https://api.amplience.net/v2/content"),
				Sensitive:   false,
			},
			"auth_url": {
				Description: "The Amplience authentication URL",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMPLIENCE_AUTH_URL", "https://auth.adis.ws/oauth/token"),
				Sensitive:   false,
			},
			"hub_id": {
				Description:      "ID of the Hub to manage",
				Type:             schema.TypeString,
				Required:         true,
				DefaultFunc:      schema.EnvDefaultFunc("AMPLIENCE_HUB_ID", nil),
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
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

	client_info := &ClientInfo{
		client: client,
		hubID:  d.Get("hub_id").(string),
	}

	return client_info, diags
}
