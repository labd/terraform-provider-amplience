package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/amplience-go-sdk/content"
	"github.com/labd/terraform-provider-amplience/amplience"
	"github.com/labd/terraform-provider-amplience/internal/resources/hub"
	"github.com/labd/terraform-provider-amplience/internal/utils"
	"net/http"
	"os"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &amplienceProvider{}
)

func New(version string) provider.Provider {
	return &amplienceProvider{
		version: version,
	}
}

type amplienceProvider struct {
	version string
}

// Provider schema struct
type amplienceProviderModel struct {
	ClientID      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
	ContentApiUrl types.String `tfsdk:"content_api_url"`
	AuthUrl       types.String `tfsdk:"auth_url"`
	HubID         types.String `tfsdk:"hub_id"`
}

// Metadata returns the provider type name.
func (p *amplienceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "amplience"
}

// Schema returns a Terraform.ResourceProvider.
func (p *amplienceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The OAuth Client ID for the Amplience management API https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication",
				Required:    true,
				Sensitive:   true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The OAuth Client Secret for Amplience management API. https://amplience_provider.com/docs/api/dynamic-content/management/index.html#section/Authentication",
				Required:    true,
				Sensitive:   true,
			},
			"content_api_url": schema.StringAttribute{
				Description: "The base URL path for the Amplience Content API",
				Optional:    true,
				Sensitive:   false,
			},
			"auth_url": schema.StringAttribute{
				Description: "The Amplience authentication URL",
				Optional:    true,
				Sensitive:   false,
			},
			"hub_id": schema.StringAttribute{
				Description: "ID of the Hub to manage",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	}
}

func (p *amplienceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config amplienceProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clientID string
	if config.ClientID.IsUnknown() || config.ClientID.IsNull() {
		clientID = os.Getenv("AMPLIENCE_CLIENT_ID")
	} else {
		clientID = config.ClientID.ValueString()
	}

	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Amplience Client ID",
			"Unknown Amplience Client ID. Please provide a valid client ID.",
		)
	}

	var clientSecret string
	if config.ClientSecret.IsUnknown() || config.ClientSecret.IsNull() {
		clientSecret = os.Getenv("AMPLIENCE_CLIENT_SECRET")
	} else {
		clientSecret = config.ClientSecret.ValueString()
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Amplience Client Secret",
			"Unknown Amplience Client Secret. Please provide a valid client secret",
		)
	}

	var hubId string
	if config.HubID.IsUnknown() || config.HubID.IsNull() {
		hubId = os.Getenv("AMPLIENCE_HUB_ID")
	} else {
		hubId = config.HubID.ValueString()
	}

	if hubId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("hub_id"),
			"Unknown Amplience Hub ID",
			"Unknown Amplience Hub ID. Please provide a valid hub ID.",
		)
	}

	var contentApiUrl string
	if config.ContentApiUrl.IsUnknown() || config.ContentApiUrl.IsNull() {
		contentApiUrl = utils.GetEnv(
			"AMPLIENCE_CONTENT_API_URL",
			"https://api.amplience.net/v2/content",
		)
	} else {
		contentApiUrl = config.ContentApiUrl.ValueString()
	}

	var authUrl string
	if config.AuthUrl.IsUnknown() || config.AuthUrl.IsNull() {
		authUrl = utils.GetEnv(
			"AMPLIENCE_AUTH_URL",
			"https://auth.amplience.net/oauth/token",
		)
	} else {
		authUrl = config.AuthUrl.ValueString()
	}

	httpClient := &http.Client{
		Transport: &utils.UserAgentTransport{
			UserAgent: fmt.Sprintf("terraform-provider-amplience/%s", p.version),
			Transport: http.DefaultTransport,
		},
	}

	client, err := content.NewClient(&content.ClientConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		URL:          contentApiUrl,
		AuthURL:      authUrl,
		HTTPClient:   httpClient,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create amplience client:\n\n"+err.Error(),
		)
		return
	}

	data := &amplience.ClientInfo{
		Client: client,
		HubID:  hubId,
	}
	resp.DataSourceData = data
	resp.ResourceData = data
}

// DataSources defines the data sources implemented in the provider.
func (p *amplienceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *amplienceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		hub.NewHubResource,
	}
}
