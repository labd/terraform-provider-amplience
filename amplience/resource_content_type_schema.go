package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentTypeSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentTypeSchemaCreate,
		ReadContext:   resourceContentTypeSchemaRead,
		UpdateContext: resourceContentTypeSchemaUpdate,
		DeleteContext: resourceContentTypeSchemaDelete,
		Schema: map[string]*schema.Schema{
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"body": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"validation_level": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceContentTypeSchemaCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id := data.Get("hub_id").(string)
	input := content.ContentTypeSchemaInput{
		SchemaID:        data.Get("schema_id").(string),
		Body:            data.Get("body").(string),
		ValidationLevel: data.Get("Validation_level").(string),
	}

	schema, err := c.ContentTypeSchemaCreate(hub_id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSchemaSaveState(data, hub_id, schema)
	return diags
}

func resourceContentTypeSchemaRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, schema_id := parseID(data.Id())

	schema, err := c.ContentTypeSchemaGet(schema_id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSchemaSaveState(data, hub_id, schema)
	return diags
}

func resourceContentTypeSchemaUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, schema_id := parseID(data.Id())

	if data.HasChange("body") || data.HasChange("validation_level") {
		current, err := c.ContentTypeSchemaGet(schema_id)
		if err != nil {
			return diag.FromErr(err)
		}

		input := content.ContentTypeSchemaInput{
			SchemaID:        data.Get("schema_id").(string),
			Body:            data.Get("body").(string),
			ValidationLevel: data.Get("Validation_level").(string),
		}

		schema, err := c.ContentTypeSchemaUpdate(current, input)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceContentTypeSchemaSaveState(data, hub_id, schema)
	}

	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentTypeSchemaDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	data.SetId("")
	return diags
}

func resourceContentTypeSchemaSaveState(data *schema.ResourceData, hub_id string, resource content.ContentTypeSchema) {
	data.SetId(createID(hub_id, resource.ID))
	data.Set("schema_id", resource.SchemaID)
	data.Set("body", resource.Body)
	data.Set("validation_level", resource.ValidationLevel)
}
