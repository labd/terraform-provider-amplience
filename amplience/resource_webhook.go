package amplience

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Description: "A webhook is a way for Dynamic Content to automatically send messages or data to a third party " +
			"system. Developers create webhooks that are triggered by specified events in Dynamic Content. These events " +
			"usually correspond to an action performed by the user such as creating or updating content, or " +
			"scheduling editions. Webhooks are associated with a single Dynamic Content hub.\n" +
			"For more info see [Amplience Webhook Docs](https://amplience.com/docs/integration/webhooks.html)",
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Description: "Label for the Webhook",
				Type:        schema.TypeString,
				Required:    true,
			},
			"events": {
				Description: "List of events to register the Webhook against",
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"handlers": {
				Description: "List of URLs to receive the Webhook",
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"active": {
				Description: "Indicates if the Webhook should be fired",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			// notifications is defined as an Array of objects in the API docs though it doesn't allow for more than
			// 1 element, throwing a "Cannot exceed the maximum of 1 notification" error if you add more so setting max
			// elements to 1
			"notifications": {
				Description: "List of notifications",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Description: "email address to notify",
							Type:        schema.TypeString,
							Required:    true,
							// TODO: Add email validation func ValidateDiagFunc:
						},
					},
				},
				MinItems: 0,
				MaxItems: 1,
			},
			"secret": {
				Description: "Shared secret between the handler and DC",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"header": {
				Description: "List of additional headers",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "Header key",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "Header value",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"secret_value": {
							Description: "Indicates whether this header value is sensitive",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Specify whether the filter is an \"in\" or an \"equal\" filter",
							Type:        schema.TypeString,
							Required:    true,
						},
						"arguments": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"json_path": {
										Description: "JSON Path of the filed you wish to match",
										Type:        schema.TypeString,
										Required:    true,
									},
									"value": {
										Description: "The value to compare too",
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
							MinItems: 0,
							MaxItems: 1,
						},
					},
				},
				MinItems: 0,
				MaxItems: 10,
			},
			"method": {
				Description: "Webhook HTTP method: POST, PATCH, PUT or DELETE",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringInSlice([]string{
					http.MethodDelete,
					http.MethodPatch,
					http.MethodPost,
					http.MethodPut,
				}, false)),
			},
			"custom_payload": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceWebhookCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	input, err := createWebhookInput(data)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating webhook draft: %w", err))
	}

	webhook, err := ci.client.WebhookCreate(ci.hubID, *input)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceWebhookSaveState(data, webhook)
	return diags
}

func resourceWebhookRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	webhook_id := data.Id()

	webhook, err := ci.client.WebhookGet(ci.hubID, webhook_id)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceWebhookSaveState(data, webhook)
	return diags
}

func resourceWebhookUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	webhook_id := data.Id()

	webhook, err := ci.client.WebhookGet(ci.hubID, webhook_id)
	if err != nil {
		return diag.FromErr(err)
	}

	input, err := createWebhookInput(data)
	if err != nil {
		return diag.FromErr(err)
	}

	new, err := ci.client.WebhookUpdate(ci.hubID, webhook, *input)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceWebhookSaveState(data, new)
	return diags
}

func resourceWebhookDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	webhook_id := data.Id()

	err := ci.client.WebhookDelete(ci.hubID, webhook_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceWebhookSaveState(data *schema.ResourceData, webhook content.Webhook) {
	data.SetId(webhook.ID)
	data.Set("label", webhook.Label)
	data.Set("events", webhook.Events)
	data.Set("handlers", webhook.Handlers)
	data.Set("active", webhook.Active)
	data.Set("secret", webhook.Secret)
	data.Set("method", webhook.Method)
	data.Set("filter", flattenWebhookFilters(&webhook.Filters))
	data.Set("custom_payload", convertCustomPayloadToMap(webhook.CustomPayload))
}

func createWebhookInput(data *schema.ResourceData) (*content.WebhookInput, error) {
	notifications := resourceWebhookGetNotifications(data.Get("notifications"))
	headers, err := resourceWebhookGetHeaders(data.Get("header"))
	if err != nil {
		return nil, fmt.Errorf("could not create webhook draft headers: %w", err)
	}

	filters, err := resourceWebhookGetFilters(data.Get("filter"))
	if err != nil {
		return nil, fmt.Errorf("could not create webhook draft filters: %w", err)
	}

	// Validation for custom payload is done in below function due to TypeMap constraints for TF provider
	customPayload, err := resourceWebhookGetCustomPayloadAndValidate(data.Get("custom_payload"))
	if err != nil {
		return nil, fmt.Errorf("error getting custom payload: %w", err)
	}

	var eventSlice []string
	for _, val := range data.Get("events").([]interface{}) {
		event := val.(string)
		if !StringInSlice([]string{
			string(content.WebhookContentItemAssigned),
			string(content.WebhookContentItemCreated),
			string(content.WebhookContentItemUpdated),
			string(content.WebhookContentItemWorkflowUpdated),
			string(content.WebhookEditionPublished),
			string(content.WebhookEditionScheduled),
			string(content.WebhookEditionUnscheduled),
			string(content.WebhookSnapshotPublished),
		}, event) {
			return nil, fmt.Errorf("invalid event type %s", event)
		}
		eventSlice = append(eventSlice, event)
	}

	var handlerSlice []string
	for _, val := range data.Get("handlers").([]interface{}) {
		handler := val.(string)
		handlerSlice = append(handlerSlice, handler)
	}

	input := &content.WebhookInput{
		Label:         data.Get("label").(string),
		Events:        eventSlice,
		Handlers:      handlerSlice,
		Active:        data.Get("active").(bool),
		Notifications: notifications,
		Secret:        data.Get("secret").(string),
		Headers:       headers,
		Filters:       filters,
		Method:        data.Get("method").(string),
		CustomPayload: customPayload,
	}
	return input, nil
}

func resourceWebhookGetNotifications(input interface{}) []content.Notification {
	inputSlice := input.([]interface{})
	var result []content.Notification

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})

		email, ok := i["email"].(string)
		if !ok {
			email = ""
		}

		result = append(result, content.Notification{Email: email})
	}

	return result
}

func resourceWebhookGetHeaders(input interface{}) ([]content.WebhookHeader, error) {
	inputSlice := input.([]interface{})
	var result []content.WebhookHeader

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})
		secret := false

		key, ok := i["key"].(string)
		if !ok {
			key = ""
		}
		value, ok := i["value"].(string)
		if !ok {
			value = ""
		}

		if secretValue, ok := i["secret_value"].(string); ok && secretValue != "" {
			value = secretValue
			secret = true
		}

		if value == "" {
			return nil, fmt.Errorf("header does not have a value defined. Specify either value or secret_value")
		}

		result = append(result, content.WebhookHeader{
			Key:    key,
			Value:  value,
			Secret: secret,
		})
	}

	return result, nil
}

func resourceWebhookGetFilters(input interface{}) ([]content.WebhookFilter, error) {
	inputSlice := input.([]interface{})
	var result []content.WebhookFilter

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})

		filterType, ok := i["type"].(string)
		if !ok {
			filterType = ""
		}
		filterArgsMap := i["arguments"].([]interface{})

		switch filterType {
		case "in":
			filter := resourceWebhookGetFilterIn(filterArgsMap)
			result = append(result, filter)
		case "equal":
			filter := resourceWebhookGetFilterEqual(filterArgsMap)
			result = append(result, filter)
		default:
			return nil, fmt.Errorf("invalid filter argument type %s", filterType)
		}
	}

	return result, nil
}

func resourceWebhookGetFilterIn(argsMap []interface{}) content.WebhookFilterIn {
	result := content.WebhookFilterIn{}

	for _, arg := range argsMap {
		j := arg.(map[string]interface{})

		if val, ok := j["json_path"].(string); ok {
			result.JSONPath = val
		}

		for _, val := range j["value"].([]interface{}) {
			if value, ok := val.(string); ok {
				result.Values = append(result.Values, value)
			}
		}
	}
	return result
}

func resourceWebhookGetFilterEqual(argsMap []interface{}) content.WebhookFilterEqual {
	result := content.WebhookFilterEqual{}

	for _, arg := range argsMap {
		j := arg.(map[string]interface{})

		if val, ok := j["json_path"].(string); ok {
			result.JSONPath = val
		}

		for _, val := range j["value"].([]interface{}) {
			if value, ok := val.(string); ok {
				result.Value = value
				break
			}
		}
	}
	return result
}

func resourceWebhookGetCustomPayloadAndValidate(input interface{}) (*content.WebhookCustomPayload, error) {
	inputMap := input.(map[string]interface{})
	payload := content.WebhookCustomPayload{}

	for key, value := range inputMap {
		if key == "type" {
			payload.Type = value.(string)
		} else if key == "value" {
			payload.Value = value.(string)
		} else {
			return nil, fmt.Errorf("unknown key %s in custom payload field", key)
		}
	}
	// If payload is empty, return nil
	if (payload == content.WebhookCustomPayload{}) {
		return nil, nil
	}

	return &payload, nil
}

func flattenWebhookFilters(filters *[]content.WebhookFilter) []interface{} {
	if filters != nil {
		fs := make([]interface{}, len(*filters))

		for i, filter := range *filters {
			f := make(map[string]interface{})

			switch v := filter.(type) {
			case content.WebhookFilterEqual:
				f["type"] = "equal"
				f["arguments"] = flattenWebhookFilterEqualArguments(v.Value, v.JSONPath)
				fs[i] = f

			case content.WebhookFilterIn:
				f["type"] = "in"
				f["arguments"] = flattenWebhookFilterInArguments(v.Values, v.JSONPath)
				fs[i] = f
			}
		}
		return fs
	}

	return make([]interface{}, 0)
}

func flattenWebhookFilterEqualArguments(Value string, JSONPath string) interface{} {
	args := make([]interface{}, 1)
	argMap := make(map[string]interface{})
	argMap["json_path"] = JSONPath
	argMap["value"] = []string{Value}
	args[0] = argMap
	return args
}

func flattenWebhookFilterInArguments(Values []string, JSONPath string) interface{} {
	args := make([]interface{}, 1)
	argMap := make(map[string]interface{})
	argMap["json_path"] = JSONPath
	argMap["value"] = Values
	args[0] = argMap
	return args
}

func convertCustomPayloadToMap(payload *content.WebhookCustomPayload) map[string]string {
	if payload != nil {
		return map[string]string{
			"type":  payload.Type,
			"value": payload.Value,
		}
	}
	return nil
}
