package amplience_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labd/terraform-provider-amplience/amplience"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ForceNew: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
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
							Required: true,
						},
						"secret": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
				ForceNew: true,
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
				ValidateDiagFunc: amplience.ValidateDiagWrapper(validation.StringInSlice([]string{
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
	c := meta.(*amplience.Client)
	APIPath := fmt.Sprintf(c.ContentAPIPath+"/hubs/%[1]s/webhooks", c.HubID)

	var webhook *amplience.Webhook
	var response *http.Response

	draft, err := createWebhookDraft(data)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating webhook draft: %w", err))
	}

	errorResponse := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		var err error
		requestBody, err := json.Marshal(draft)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error marshalling draft %v: %w", draft, err))
		}
		response, err = amplience.AmplienceRequest(APIPath, http.MethodPost, bytes.NewBuffer(requestBody))
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error during http request: %w", err))
		}
		return amplience.HandleAmplienceError(response)
	})

	if errorResponse != nil {
		return diag.FromErr(fmt.Errorf("received error from request, could not create webhook for draft %v", draft))
	}

	if response == nil {
		return diag.FromErr(fmt.Errorf("could not create webhook"))
	}
	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &webhook)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing response body into webhook struct: %w", err))
	}
	data.SetId(webhook.ID)

	resourceWebhookRead(ctx, data, meta)

	return diags
}

func resourceWebhookRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	webhookID := data.Id()
	// Below to be replaced with client library function
	webhook, err := getWebhookWithID(webhookID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred when trying to get webhook with id %[1]s",
			webhookID))
	}
	if webhook == nil {
		log.Print("[DEBUG] No webhook found")
		data.SetId("")
	} else {
		log.Print("[DEBUG] Found following webhook: ")
		log.Print(amplience.StringFormatObject(webhook))

		data.Set("label", webhook.Label)
		data.Set("events", webhook.Events)
		data.Set("handlers", webhook.Handlers)
		data.Set("active", webhook.Active)

		data.Set("secret", webhook.Secret)
		filters := flattenWebhookFilters(&webhook.Filters)
		data.Set("filter", filters)
		data.Set("method", webhook.Method)
		data.Set("custom_payload", convertCustomPayloadToMap(webhook.CustomPayload))

		// NOTE: We don't set 'headers' and 'notifications' here as their response can come back as nulls leading to a
		// state difference. In order to avoid any mismatching state issues we set ForceNew to true for both fields
		// so a new resource is created if there are changes in either field
	}
	return diags
}

func resourceWebhookUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	webhookID := data.Id()

	// Ideally we'd be able to specify a specific PATCH request which only contains the fields that have data.HasChange
	// this can probably best be abstracted into the SDK
	draft, err := createWebhookDraft(data)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating webhook draft: %w", err))
	}
	// We force the creation of a new resource upon any change to headers or notifications due to their secret behaviour
	// This function should not be called when there are changes to these fields but drop them from the draft to make
	//doubly sure we avoid touching them in the PATCH request
	draft.Headers = nil
	draft.Notifications = nil

	requestBody, err := json.Marshal(draft)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not marshal %v", draft))
	}

	webhook, err := updateWebhookWithID(webhookID, bytes.NewBuffer(requestBody), meta)
	if webhook == nil {
		log.Printf("[DEBUG] Nothing update for webhook with ID: %s", webhookID)
	} else {
		log.Printf("[DEBUG] Succesfully updated webhook with ID: %s", webhookID)
		data.Set("label", webhook.Label)
		data.Set("events", webhook.Events)
		data.Set("handlers", webhook.Handlers)
		data.Set("active", webhook.Active)

		data.Set("secret", webhook.Secret)
		filters := flattenWebhookFilters(&webhook.Filters)
		data.Set("filter", filters)
		data.Set("method", webhook.Method)
		data.Set("custom_payload", convertCustomPayloadToMap(webhook.CustomPayload))
		// NOTE: We don't set 'headers' and 'notifications' here as their response can come back as nulls leading to a
		// state difference. In order to avoid any mismatching state issues we set ForceNew to true for both fields
		// so a new resource is created if there are changes in either field
	}

	resourceWebhookRead(ctx, data, meta)

	return diags
}

func resourceWebhookDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	webhookID := data.Id()

	c := meta.(*amplience.Client)

	webhook, err := getWebhookWithID(webhookID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not get Webhook with ID %s for Hub %s", webhookID, c.HubID))
	}
	if webhook == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprint("No Webhook found, nothing to delete"),
		})
		return diags
	}

	err = deleteWebhookWithID(webhookID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete Webhook with ID %s", webhookID))
	}
	return diags
}

func createWebhookDraft(data *schema.ResourceData) (*amplience.Webhook, error) {
	notifications := resourceWebhookGetNotifications(data.Get("notifications"))
	headers := resourceWebhookGetHeaders(data.Get("header"))
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
		if !amplience.StringInSlice([]string{
			string(amplience.WebhookContentItemAssigned),
			string(amplience.WebhookContentItemCreated),
			string(amplience.WebhookContentItemUpdated),
			string(amplience.WebhookContentItemWorkflowUpdated),
			string(amplience.WebhookEditionPublished),
			string(amplience.WebhookEditionScheduled),
			string(amplience.WebhookEditionUnscheduled),
			string(amplience.WebhookSnapshotPublished),
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

	draft := &amplience.Webhook{
		Label:         data.Get("label").(string),
		Events:        eventSlice,
		Handlers:      handlerSlice,
		Active:        data.Get("active").(bool),
		Notifications: notifications,
		Secret:        data.Get("secret").(string),
		Headers:       headers,
		Filters:       filters,
		Method:        data.Get("method").(string),
		CustomPayload: *customPayload,
	}
	return draft, nil
}

// getWebhookWithID returns a webhook based on a hubID and a webhookID if it exists. Else returns nil.
// The functionality of this should be abstracted into client library
func getWebhookWithID(webhookID string, meta interface{}) (*amplience.Webhook, error) {
	webhook := amplience.Webhook{}

	c := meta.(*amplience.Client)
	APIPath := fmt.Sprintf(c.ContentAPIPath+"/hubs/%[1]s/webhooks/%[2]s", c.HubID, webhookID)

	response, err := amplience.AmplienceRequest(APIPath, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make GET request to %s: %w", APIPath, err)
	}
	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &webhook)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into webhook struct: %w", err)
	}
	return &webhook, nil
}

// deleteWebhookWithID deletes a webhook based on a hubID and a webhookID if it exists.
// The functionality of this should be abstracted into client library
// Note that in its current state the function can not detect whether a resource is present before deleting it
func deleteWebhookWithID(webhookID string, meta interface{}) error {
	c := meta.(*amplience.Client)

	APIPath := fmt.Sprintf(c.ContentAPIPath+"/hubs/%[1]s/webhooks/%[2]s", c.HubID, webhookID)

	response, err := amplience.AmplienceRequest(APIPath, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("unable to make DELETE request to %s: %w", APIPath, err)
	}
	if response.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("received unexpected status code %d", response.StatusCode)
}

// updateWebhookWithID sends a PATCH request to update a Webhook based on a HubID and a WebhookID
// The functionality of this should be abstracted into client library
func updateWebhookWithID(webhookID string, requestBody *bytes.Buffer, meta interface{}) (*amplience.Webhook, error) {
	webhook := amplience.Webhook{}
	c := meta.(*amplience.Client)

	APIPath := fmt.Sprintf(c.ContentAPIPath+"/hubs/%[1]s/webhooks/%[2]s", c.HubID, webhookID)

	response, err := amplience.AmplienceRequest(APIPath, http.MethodPatch, requestBody)
	if err != nil {
		return nil, fmt.Errorf("unable to make PATCH request to %s: %w", APIPath, err)
	}

	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &webhook)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into webhook struct: %w", err)
	}

	return &webhook, nil
}

func resourceWebhookGetNotifications(input interface{}) []amplience.Notification {
	inputSlice := input.([]interface{})
	var result []amplience.Notification

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})

		email, ok := i["email"].(string)
		if !ok {
			email = ""
		}

		result = append(result, amplience.Notification{Email: email})
	}

	return result
}

func resourceWebhookGetHeaders(input interface{}) []amplience.WebhookHeader {
	inputSlice := input.([]interface{})
	var result []amplience.WebhookHeader

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})

		key, ok := i["key"].(string)
		if !ok {
			key = ""
		}
		value, ok := i["value"].(string)
		if !ok {
			value = ""
		}
		secret, ok := i["secret"].(bool)
		if !ok {
			secret = false
		}

		result = append(result, amplience.WebhookHeader{
			Key:    key,
			Value:  value,
			Secret: secret,
		})
	}

	return result
}

func resourceWebhookGetFilters(input interface{}) ([]amplience.WebhookFilter, error) {
	inputSlice := input.([]interface{})
	var result []amplience.WebhookFilter

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})
		var singleVal bool

		filterType, ok := i["type"].(string)
		if !ok {
			filterType = ""
		}
		switch filterType {
		case "in":
			singleVal = false
		case "equal":
			singleVal = true
		default:
			return nil, fmt.Errorf("invalid filter argument type %s", filterType)
		}

		var filterArgs []amplience.RawArg
		filterArgsMap, ok := i["arguments"].([]interface{})
		if !ok {
			filterArgs = nil
		} else {
			filterArgs = resourceWebhookGetFilterArgs(filterArgsMap, singleVal)
		}

		result = append(result, amplience.WebhookFilter{
			Type:      filterType,
			Arguments: filterArgs,
		})
	}

	return result, nil
}

func resourceWebhookGetFilterArgs(filterArgsMap []interface{}, singleVal bool) []amplience.RawArg {
	if singleVal {
		return resourceWebhookGetEqualFilterArgs(filterArgsMap)
	} else {
		return resourceWebhookGetInFilterArgs(filterArgsMap)
	}
}

func resourceWebhookGetInFilterArgs(argsMap []interface{}) []amplience.RawArg {
	var values []string
	for _, arg := range argsMap {
		j := arg.(map[string]interface{})

		jsonPath, ok := j["json_path"].(string)
		if !ok {
			jsonPath = ""
		}
		for _, val := range j["value"].([]interface{}) {
			value, ok := val.(string)
			if !ok {
				value = ""
			}
			values = append(values, value)
		}
		return []amplience.RawArg{{
			JSONPath: &jsonPath,
		}, {
			InValues: &values,
		}}
	}
	return nil
}

func resourceWebhookGetEqualFilterArgs(argsMap []interface{}) []amplience.RawArg {
	for _, arg := range argsMap {
		j := arg.(map[string]interface{})

		jsonPath, ok := j["json_path"].(string)
		if !ok {
			jsonPath = ""
		}

		for _, val := range j["value"].([]interface{}) {
			value, ok := val.(string)
			if !ok {
				value = ""
			}
			return []amplience.RawArg{{
				JSONPath: &jsonPath,
			}, {
				EqValue: &value,
			}}
		}

	}
	return nil
}

func resourceWebhookGetCustomPayloadAndValidate(input interface{}) (*amplience.WebhookCustomPayload, error) {
	inputMap := input.(map[string]interface{})
	payload := amplience.WebhookCustomPayload{}

	for key, value := range inputMap {
		if key == "type" {
			payload.Type = value.(string)
		} else if key == "value" {
			payload.Value = value.(string)
		} else {
			return nil, fmt.Errorf("unknown key %s in custom payload field", key)
		}
	}

	return &payload, nil
}

func flattenWebhookFilters(filters *[]amplience.WebhookFilter) []interface{} {
	if filters != nil {
		fs := make([]interface{}, len(*filters), len(*filters))

		for i, filter := range *filters {
			f := make(map[string]interface{})

			f["type"] = filter.Type
			f["arguments"] = flattenWebhookFilterArguments(filter.Arguments, filter.Type)
			fs[i] = f
		}

		return fs
	}
	return make([]interface{}, 0)
}

func flattenWebhookFilterArguments(arguments []amplience.RawArg, filterType string) interface{} {
	// We know its a list of 2 elements of which the first has jsonPath and the second has its value so....
	// TODO: cleam me up and everyrthign else
	args := make([]interface{}, 1, 1)
	argMap := make(map[string]interface{})
	jsonPath := arguments[0].JSONPath
	argMap["json_path"] = jsonPath
	if filterType == "equal" {
		eqValueSlice := make([]string, 1)
		eqValueSlice[0] = *arguments[1].EqValue
		argMap["value"] = eqValueSlice
	} else if filterType == "in" {
		argMap["value"] = arguments[1].InValues
	}
	args[0] = argMap
	return args
}

func convertCustomPayloadToMap(payload amplience.WebhookCustomPayload) map[string]string {
	return map[string]string{
		"type":  payload.Type,
		"value": payload.Value,
	}
}
