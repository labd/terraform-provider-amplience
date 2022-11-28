package amplience

import (
	"context"
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
			"auto_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	instance, err := ci.client.ContentTypeSchemaCreate(ci.hubID, input)

	if errResp, ok := err.(*content.ErrorResponse); ok {
		if errResp.StatusCode >= 400 {

			log.Println("Received 400 conflict response: content type schema already exists.")
			log.Println("Proceeding to unarchive if necessary and update exiting content type schema.")

			instance, err = ci.client.ContentTypeSchemaFindBySchemaId(input.SchemaID, ci.hubID)
			if err != nil {
				return diag.FromErr(err)
			}

			if instance.Status == string(content.StatusArchived) {
				instance, err = ci.client.ContentTypeSchemaUnarchive(instance.ID, instance.Version)
			}
			if err != nil {
				return diag.FromErr(err)
			}

			instance, err = ci.client.ContentTypeSchemaUpdate(instance, input)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSchemaSaveState(data, instance)
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

	id := data.Id()

	if data.HasChange("body") || data.HasChange("validation_level") {
		instance, err := ci.client.ContentTypeSchemaGet(id)

		if err != nil {
			return diag.FromErr(err)
		}

		if instance.Status == string(content.StatusArchived) {
			log.Println("Content type schema was archived. Proceed to unarchive first before applying update.")
			instance, err = ci.client.ContentTypeSchemaUnarchive(instance.ID, instance.Version)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		input := content.ContentTypeSchemaInput{
			SchemaID:        data.Get("schema_id").(string),
			Body:            data.Get("body").(string),
			ValidationLevel: data.Get("validation_level").(string),
		}
		schema, err := ci.client.ContentTypeSchemaUpdate(instance, input)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceContentTypeSchemaSaveState(data, schema)

		if data.Get("auto_sync").(bool) {
			log.Printf("Content type schema %s has been updated; will auto-sync...", instance.SchemaID)

			contentType, err := ci.client.ContentTypeFindByUri(instance.SchemaID, ci.hubID)
			if err != nil {
				log.Printf("No content type found for %s, will skip syncing", instance.SchemaID)
			}

			syncResult, err := ci.client.ContentTypeSyncSchema(contentType)

			if err != nil {
				// Log, but don't fail the apply step itself
				log.Fatalf("Could not sync content type: %s", err)
			}

			log.Printf("Synced content type %s", syncResult.ContentTypeURI)
		}
	}

	return diags
}

// The Amplience API does not have a content type delete functionality. Setting ID to "" and returning nil
func resourceContentTypeSchemaDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()
	instance, err := ci.client.ContentTypeSchemaGet(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if instance.Status == string(content.StatusActive) {
		_, err = ci.client.ContentTypeSchemaArchive(id, instance.Version)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	data.Set("schema_id", instance.SchemaID)
	data.Set("body", instance.Body)
	data.Set("validation_level", instance.ValidationLevel)
	data.Set("version", instance.Version)
	return diags
}

func resourceContentTypeSchemaSaveState(data *schema.ResourceData, resource content.ContentTypeSchema) {
	data.SetId(resource.ID)
	data.Set("schema_id", resource.SchemaID)
	data.Set("body", resource.Body)
	data.Set("validation_level", resource.ValidationLevel)
	data.Set("version", resource.Version)
}
