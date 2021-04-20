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
		Description: "Content type schemas are JSON schemas that define a type of content to be created, including its " +
			"structure, format and validation rules. In Dynamic Content, content type schemas match the format of the " +
			"JSON Schema standard, with a few extensions and some keywords that are not supported.\n" +
			"For more info see [Amplience Content Type Schema Docs](https://amplience.com/docs/integration/contenttypes.html)",
		CreateContext: resourceContentTypeSchemaCreate,
		ReadContext:   resourceContentTypeSchemaRead,
		UpdateContext: resourceContentTypeSchemaUpdate,
		DeleteContext: resourceContentTypeSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
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
		ValidationLevel: data.Get("validation_level").(string),
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

	schema_id := data.Id()
	hub_id := data.Get("hub_id").(string)

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

	schema_id := data.Id()
	hub_id := data.Get("hub_id").(string)
	if data.HasChange("body") || data.HasChange("validation_level") {
		current, err := c.ContentTypeSchemaGet(schema_id)
		if err != nil {
			return diag.FromErr(err)
		}

		input := content.ContentTypeSchemaInput{
			SchemaID:        data.Get("schema_id").(string),
			Body:            data.Get("body").(string),
			ValidationLevel: data.Get("validation_level").(string),
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
	c := m.(*content.Client)

	id := data.Id()
	version := data.Get("version").(int)

	_, err := c.ContentTypeSchemaArchive(id, version)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceContentTypeSchemaSaveState(data *schema.ResourceData, hub_id string, resource content.ContentTypeSchema) {
	data.SetId(resource.ID)
	data.Set("hub_id", hub_id)
	data.Set("schema_id", resource.SchemaID)
	data.Set("body", resource.Body)
	data.Set("validation_level", resource.ValidationLevel)
	data.Set("version", resource.Version)
}
