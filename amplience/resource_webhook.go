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
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
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
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"events": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"handlers": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// notifications is defined as an Array of objects in the API docs though it doesn't allow for more than
			// 1 element, throwing a "Cannot exceed the maximum of 1 notification" error if you add more so setting max
			// elements to 1
			"notifications": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
							// TODO: Add email validation func ValidateDiagFunc:
						},
					},
				},
				MinItems: 0,
				MaxItems: 1,
			},
			"secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"header": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"arguments": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"json_path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
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
				Type:     schema.TypeString,
				Required: true,
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
	c := meta.(*content.Client)

	hub_id := data.Get("hub_id").(string)

	input, err := createWebhookInput(data)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating webhook draft: %w", err))
	}

	webhook, err := c.WebhookCreate(hub_id, *input)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceWebhookSaveState(data, hub_id, webhook)
	return diags
}

func resourceWebhookRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, webhook_id := parseID(data.Id())

	webhook, err := c.WebhookGet(hub_id, webhook_id)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceWebhookSaveState(data, hub_id, webhook)
	return diags
}

func resourceWebhookUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, webhook_id := parseID(data.Id())

	webhook, err := c.WebhookGet(hub_id, webhook_id)
	if err != nil {
		return diag.FromErr(err)
	}

	input, err := createWebhookInput(data)
	if err != nil {
		return diag.FromErr(err)
	}

	new, err := c.WebhookUpdate(hub_id, webhook, *input)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceWebhookSaveState(data, hub_id, new)
	return diags
}

func resourceWebhookDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id, webhook_id := parseID(data.Id())

	err := c.WebhookDelete(hub_id, webhook_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceWebhookSaveState(data *schema.ResourceData, hub_id string, webhook content.Webhook) {
	data.SetId(createID(hub_id, webhook.ID))
	data.Set("label", webhook.Label)
	data.Set("events", webhook.Events)
	data.Set("handlers", webhook.Handlers)
	data.Set("active", webhook.Active)
	data.Set("secret", webhook.Secret)
	data.Set("method", webhook.Method)
	data.Set("filter", flattenWebhookFilters(webhook.Filters))
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
		filterArgsMap, ok := i["arguments"].([]interface{})

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

func flattenWebhookFilters(filters []content.WebhookFilter) []interface{} {
	result := make([]interface{}, len(filters))

	for _, filter := range filters {
		item := make(map[string]interface{})
		arguments := make(map[string]interface{})

		switch v := filter.(type) {
		case content.WebhookFilterEqual:
			item["type"] = "equal"
			arguments["json_path"] = v.JSONPath
			arguments["values"] = []string{v.Value}
			item["arguments"] = arguments

		case content.WebhookFilterIn:
			item["type"] = "in"
			arguments["json_path"] = v.JSONPath
			arguments["values"] = v.Values
			item["arguments"] = arguments
		}

	}

	return result
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
