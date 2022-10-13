package amplience

import (
	"context"
	"fmt"
	"log"

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
	ci := getClient(meta)
	schemaId := data.Get("schema_id").(string)

	input := content.ContentTypeSchemaInput{
		SchemaID:        schemaId,
		Body:            data.Get("body").(string),
		ValidationLevel: data.Get("validation_level").(string),
	}

	schema, err := ci.client.ContentTypeSchemaCreate(ci.hubID, input)

	if errResp, ok := err.(*content.ErrorResponse); ok {
		if errResp.StatusCode == 409 {
			log.Println("Received 409 conflict response; schema must be unarchived")
			schema, err = _unarchiveSchema(schemaId, ci)
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSchemaSaveState(data, schema)
	return diags
}

func resourceContentTypeSchemaRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	schema_id := data.Id()
	schema, err := ci.client.ContentTypeSchemaGet(schema_id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSchemaSaveState(data, schema)
	return diags
}

func resourceContentTypeSchemaUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	schema_id := data.Id()
	if data.HasChange("body") || data.HasChange("validation_level") {
		current, err := ci.client.ContentTypeSchemaGet(schema_id)
		if err != nil {
			return diag.FromErr(err)
		}

		input := content.ContentTypeSchemaInput{
			SchemaID:        data.Get("schema_id").(string),
			Body:            data.Get("body").(string),
			ValidationLevel: data.Get("validation_level").(string),
		}

		schema, err := ci.client.ContentTypeSchemaUpdate(current, input)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceContentTypeSchemaSaveState(data, schema)
	}

	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentTypeSchemaDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()
	version := data.Get("version").(int)

	_, err := ci.client.ContentTypeSchemaArchive(id, version)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceContentTypeSchemaSaveState(data *schema.ResourceData, resource content.ContentTypeSchema) {
	data.SetId(resource.ID)
	data.Set("schema_id", resource.SchemaID)
	data.Set("body", resource.Body)
	data.Set("validation_level", resource.ValidationLevel)
	data.Set("version", resource.Version)
}

func _unarchiveSchema(schemaId string, ci *ClientInfo) (content.ContentTypeSchema, error) {
	result := content.ContentTypeSchema{}

	log.Printf("Get info for content type schema %s", schemaId)
	schema, getErr := _schemaBySchemaId(schemaId, ci, content.StatusAny)

	if getErr != nil {
		return result, getErr
	}

	log.Printf("Received content type schema version %v", schema.Version)
	schema, err := ci.client.ContentTypeSchemaUnarchive(schema.ID, schema.Version)

	if err != nil {
		return schema, err
	}

	return schema, err
}

func _schemaBySchemaId(schemaId string, ci *ClientInfo, status content.ContentStatus) (content.ContentTypeSchema, error) {
	dummy := content.ContentTypeSchema{}
	schemaList, getErr := ci.client.ContentTypeSchemaGetAll(ci.hubID, status)

	if getErr != nil {
		return dummy, getErr
	}

	for _, schema := range schemaList {
		if schema.SchemaID == schemaId {
			return schema, nil
		}
	}

	return dummy, fmt.Errorf(fmt.Sprintf("Could not find schema %s", schemaId))
}
