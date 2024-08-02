package hub

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/labd/amplience-go-sdk/content"
	"github.com/labd/terraform-provider-amplience/amplience"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &hubResource{}
	_ resource.ResourceWithConfigure   = &hubResource{}
	_ resource.ResourceWithImportState = &hubResource{}
)

// NewHubResource is a helper function to simplify the provider implementation.
func NewHubResource() resource.Resource {
	return &hubResource{}
}

// hubResource is the resource implementation.
type hubResource struct {
	client *content.Client
	hubId  string
}

// Metadata returns the data source type name.
func (r *hubResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hub"
}

// Schema defines the schema for the data source.
func (r *hubResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Permissions are set at the hub level. All users of a hub can at least view all of the " +
			"content within the repositories inside that hub. Content cannot be shared across hubs. However, content " +
			"can be shared and linked to across repositories within the same hub. So you can create a content item " +
			"in one repository and include content stored in another. Events and editions are scheduled within a " +
			"single hub. So if you want an overall view of the planning calendar across many brands, then you may wish " +
			"to consider a single hub. However, in some cases you may want to keep the calendars separate. Many " +
			"settings, such as the publishing endpoint (the subdomain to which your content is published) are set at " +
			"a hub level. Multiple hubs may publish content to the same endpoint.\n\n" +
			"For more info see [Amplience Hubs & Repositories Docs](https://amplience.com/docs/intro/hubsandrepositories.html)\n\n" +
			"**It is recommended to import a " +
			"new hub instead of creating it!** This is because the hub already exists, so any differences in " +
			"configuration might be overwritten, leading to unintended outcomes.",
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "Hub name",
				Required:    true,
			},
			"label": schema.StringAttribute{
				Description: "Hub label",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Hub description",
				Optional:    true,
			},
			"settings": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Hub settings",
				Attributes: map[string]schema.Attribute{
					"publishing": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"platforms": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"amplience_dam": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"api_key": schema.StringAttribute{
												Description: "DAM publishing client key",
												Required:    true,
											},
											"api_secret": schema.StringAttribute{
												Description: "DAM publishing client secret",
												Required:    true,
												Sensitive:   true,
											},
											"endpoint": schema.StringAttribute{
												Description: "Publishing endpoint, also known as Company Tag",
												Required:    true,
											},
										},
									},
								},
							},
						},
					},
					"devices": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:   true,
									Validators: []validator.String{stringvalidator.LengthBetween(1, 50)},
								},
								"width": schema.Int64Attribute{
									Required: true,
								},
								"height": schema.Int64Attribute{
									Required: true,
								},
								"orientate": schema.BoolAttribute{
									Required: true,
								},
							},
						},
					},
					"localization": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"locales": schema.ListAttribute{
								Optional:    true,
								ElementType: basetypes.StringType{},
							},
						},
					},
					"applications": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"template_uri": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"preview_virtual_staging_environment": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"hostname": schema.StringAttribute{
								Description: "Virtual Staging Environment hostname",
								Required:    true,
							},
						},
					},
					"virtual_staging_environment": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"hostname": schema.StringAttribute{
								Description: "Virtual Staging Environment hostname",
								Required:    true,
							},
						},
					},
					"asset_management": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional: true,
							},
							"client_config": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("HUB", "USER"),
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *hubResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data := req.ProviderData.(*amplience.ClientInfo)
	r.client = data.Client
	r.hubId = data.HubID
}

// Create creates the resource and sets the initial Terraform state.
func (r *hubResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Hub
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hub, err := r.client.HubGet(r.hubId)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get hub", err.Error())
		return
	}
	current := NewHubFromNative(&hub)
	fmt.Println("current", current)

	hub, err = r.client.HubPatch(r.hubId, current.ToUpdateInput())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update hub", err.Error())
		return
	}

	result := NewHubFromNative(&hub)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *hubResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Hub
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.HubGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading hub", err.Error())
		return
	}
	current := NewHubFromNative(&res)
	current.setSecretValuesFromState(state)

	// Set refreshed state
	diags = resp.State.Set(ctx, current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *hubResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var current Hub
	diags := req.State.Get(ctx, &current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get updated plan
	var plan Hub
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	hub, err := r.client.HubPatch(current.ID.ValueString(), plan.ToUpdateInput())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update hub", err.Error())
		return
	}

	newState := NewHubFromNative(&hub)
	newState.setSecretValuesFromState(current)

	// Set updated state
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *hubResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Delete not implemented", "Deleting a hub is not supported. The hub data has been removed from state, but still exists in Amplience.")
}

func (r *hubResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
