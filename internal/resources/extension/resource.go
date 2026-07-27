package extension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/labd/amplience-go-sdk/content"
	"github.com/labd/terraform-provider-amplience/amplience"
)

var (
	_ resource.Resource                = &extensionResource{}
	_ resource.ResourceWithConfigure   = &extensionResource{}
	_ resource.ResourceWithImportState = &extensionResource{}
)

func NewExtensionResource() resource.Resource {
	return &extensionResource{}
}

type extensionResource struct {
	client *content.Client
	hubId  string
}

func (r *extensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extension"
}

func (r *extensionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Amplience Extension. Extensions are UI components embedded in Dynamic Content to provide custom editing experiences.\n\n" +
			"For more info see [Amplience Extensions Docs](https://amplience.com/docs/development/extensions.html)",
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Extension name",
				Required:    true,
			},
			"label": schema.StringAttribute{
				Description: "Extension label displayed in the UI",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Extension description",
				Optional:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL of the extension",
				Required:    true,
			},
			"height": schema.Int64Attribute{
				Description: "Height of the extension iframe in pixels. Only applicable for CONTENT_FIELD category",
				Optional:    true,
				Computed:    true,
			},
			"enabled_for_all_content_types": schema.BoolAttribute{
				Description: "Whether this extension is enabled for all content types. Only applicable for CONTENT_EDITOR category",
				Optional:    true,
				Computed:    true,
			},
			"category": schema.StringAttribute{
				Description: "Extension category (CONTENT_FIELD, DASHBOARD, or CONTENT_EDITOR). Determines which other attributes are applicable",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("CONTENT_FIELD", "DASHBOARD", "CONTENT_EDITOR"),
				},
			},
			"parameters": schema.StringAttribute{
				Description: "Extension parameters as a JSON object string, e.g. jsonencode({vendor = \"commercetools\"})",
				Optional:    true,
			},
			"snippets": schema.ListNestedAttribute{
				Description: "Extension snippets. Only applicable for CONTENT_FIELD category",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label": schema.StringAttribute{
							Description: "Snippet label",
							Required:    true,
						},
						"body": schema.StringAttribute{
							Description: "Snippet body as a JSON string",
							Required:    true,
						},
					},
				},
			},
			"settings": schema.StringAttribute{
				Description: "Extension settings as a JSON string",
				Optional:    true,
			},
		},
	}
}

func (r *extensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data := req.ProviderData.(*amplience.ClientInfo)
	r.client = data.Client
	r.hubId = data.HubID
}

func (r *extensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Extension
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := plan.ToInput(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ext, err := r.client.ExtensionCreate(r.hubId, input)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create extension", err.Error())
		return
	}

	state := NewExtensionFromNative(&ext, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *extensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Extension
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ext, err := r.client.ExtensionGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read extension", err.Error())
		return
	}

	newState := NewExtensionFromNative(&ext, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (r *extensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Extension
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan Extension
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch current state from API to get the ID for the update endpoint
	current, err := r.client.ExtensionGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read extension for update", err.Error())
		return
	}

	input := plan.ToInput(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ext, err := r.client.ExtensionUpdate(current, input)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update extension", err.Error())
		return
	}

	newState := NewExtensionFromNative(&ext, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (r *extensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Extension
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ExtensionDelete(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Unable to delete extension", err.Error())
	}
}

func (r *extensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
