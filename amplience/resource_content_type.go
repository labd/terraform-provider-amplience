package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentTypeCreate,
		ReadContext:   resourceContentTypeRead,
		UpdateContext: resourceContentTypeUpdate,
		DeleteContext: resourceContentTypeDelete,
		Schema: map[string]*schema.Schema{
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"content_type_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceContentTypeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id := data.Get("hub_id").(string)
	input := resourceContentTypeCreateInput(data)
	content_type, err := c.ContentTypeCreate(hub_id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, hub_id, content_type)
	return diags
}

func resourceContentTypeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, content_type_id := parseID(data.Id())

	content_type, err := c.ContentTypeGet(content_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, hub_id, content_type)
	return diags
}

func resourceContentTypeUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, content_type_id := parseID(data.Id())

	if data.HasChange("body") || data.HasChange("validation_level") {
		current, err := c.ContentTypeGet(content_type_id)
		if err != nil {
			return diag.FromErr(err)
		}

		input := resourceContentTypeCreateInput(data)
		content_type, err := c.ContentTypeUpdate(current, input)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceContentTypeSaveState(data, hub_id, content_type)
	}

	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentTypeDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	data.SetId("")
	return diags
}

func resourceContentTypeSaveState(data *schema.ResourceData, hub_id string, resource content.ContentType) {
	data.SetId(createID(hub_id, resource.ID))
	data.Set("content_type_uri", resource.ContentTypeURI)
	data.Set("status", resource.Status)
	data.Set("label", resource.Settings.Label)
}

func resourceContentTypeCreateInput(data *schema.ResourceData) content.ContentTypeInput {
	return content.ContentTypeInput{
		ContentTypeURI: data.Get("content_type_uri").(string),
		Settings: content.ContentTypeSettings{
			Label: data.Get("label").(string),
		},
	}
}
