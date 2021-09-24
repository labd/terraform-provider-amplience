package amplience

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/labd/amplience-go-sdk/content"
)

func resourceSearchIndex() *schema.Resource {
	return &schema.Resource{
		Description: "A search index is the connection between Amplience and Algolia." +
			"For more info see [Amplience Index Docs](https://amplience.com/docs/development/search-indexes/readme.html)",
		CreateContext: resourceSearchIndexCreate,
		ReadContext:   resourceSearchIndexRead,
		UpdateContext: resourceSearchIndexUpdate,
		DeleteContext: resourceSearchIndexDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Description: "Label for the Index",
				Type:        schema.TypeString,
				Required:    true,
			},
			"suffix": {
				Description: "Suffix for the Index",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "Either PRODUCTION or STAGING",
				Type:        schema.TypeString,
				Required:    true,
			},
			"content_types": {
				Description: "List of content type urls. Each content type will create 2 corresponding webhooks (PUT & DELETE)",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"settings": {
				Type:        schema.TypeString,
				Description: "A JSON string containing Algolia settings (https://www.algolia.com/doc/api-reference/api-parameters/)",
				Optional:    true,
			},
			"webhook_custom_payload": {
				Type:        schema.TypeMap,
				Description: "A Handlebars Json string for the custom payload that will be used for each content type webhook",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSearchIndexCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	input, err := createIndexInput(data)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating index draft: %w", err))
	}

	resource, err := ci.client.AlgoliaIndexCreate(ci.hubID, *input)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateIndexWebhooksAndSettings(ci.client, ci.hubID, resource.ID, data)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSearchIndexSaveState(data, resource)
	return diags
}

func resourceSearchIndexRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()

	resource, err := ci.client.AlgoliaIndexGet(ci.hubID, id)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceSearchIndexSaveState(data, resource)
	return diags
}

func resourceSearchIndexUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()

	old, err := ci.client.AlgoliaIndexGet(ci.hubID, id)
	if err != nil {
		return diag.FromErr(err)
	}

	input, err := createIndexInput(data)
	if err != nil {
		return diag.FromErr(err)
	}

	new, err := ci.client.AlgoliaIndexUpdate(ci.hubID, old, *input)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateIndexWebhooksAndSettings(ci.client, ci.hubID, id, data)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSearchIndexSaveState(data, new)
	return diags
}

func resourceSearchIndexDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()

	_, err := ci.client.AlgoliaIndexDelete(ci.hubID, id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceSearchIndexSaveState(data *schema.ResourceData, resource content.AlgoliaIndex) {
	data.SetId(resource.ID)
	data.Set("label", resource.Label)
	data.Set("suffix", resource.Suffix)
	data.Set("type", resource.Type)
}

func createIndexInput(data *schema.ResourceData) (*content.AlgoliaIndexInput, error) {

	var assignedContentTypes []content.AssignedContentTypeInput
	for _, val := range data.Get("content_types").([]interface{}) {
		uri := val.(string)
		assignedContentTypes = append(assignedContentTypes, content.AssignedContentTypeInput{
			ContentTypeUri: uri,
		})
	}

	indexType := data.Get("type").(string)
	if indexType != "PRODUCTION" && indexType != "STAGING" {
		return nil, fmt.Errorf("type must be either 'PRODUCTION' or 'STAGING'")
	}

	var handlerSlice []string
	for _, val := range data.Get("handlers").([]interface{}) {
		handler := val.(string)
		handlerSlice = append(handlerSlice, handler)
	}

	input := &content.AlgoliaIndexInput{
		Label:                data.Get("label").(string),
		Suffix:               data.Get("suffix").(string),
		Type:                 data.Get("type").(string),
		AssignedContentTypes: assignedContentTypes,
	}
	return input, nil
}

func createAlgoliaIndexSettings(input string) (content.AlgoliaIndexSettings, error) {
	settingsInput := content.AlgoliaIndexSettings{}
	err := json.Unmarshal([]byte(input), &settingsInput)
	return settingsInput, err
}
func updateIndexWebhooksAndSettings(client *content.Client, hubId string, indexId string, data *schema.ResourceData) error {

	settingsInput, err := createAlgoliaIndexSettings(data.Get("settings").(string))
	if err != nil {
		return err
	}

	_, err = client.AlgoliaIndexSettingsUpdate(hubId, indexId, settingsInput)
	if err != nil {
		return err
	}

	customPayload, err := resourceWebhookGetCustomPayloadAndValidate(data.Get("webhook_custom_payload"))
	if err != nil {
		return err
	}

	webhooks, err := client.AlgoliaIndexWebhooksGet(hubId, indexId)
	if err != nil {
		return err
	}

	for _, item := range webhooks {
		_, err = client.WebhookUpdate(hubId, item, content.WebhookInput{
			CustomPayload: customPayload,
			Label:         item.Label,
			Events:        item.Events,
			Active:        item.Active,
			Notifications: item.Notifications,
			Secret:        item.Secret,
			Filters:       item.Filters,
			Method:        item.Method,
			Handlers:      item.Handlers,
		})
		if err != nil {
			return err
		}
	}

	return err
}
