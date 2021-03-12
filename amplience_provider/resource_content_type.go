package amplience_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/labd/terraform-provider-amplience/amplience"
)

func resourceContentType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentTypeCreate,
		ReadContext:   resourceContentTypeRead,
		UpdateContext: resourceContentTypeUpdate,
		DeleteContext: resourceContentTypeDelete,
		Schema: map[string]*schema.Schema{
			"content_type_uri": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: amplience.ValidateDiagWrapper(validation.IsURLWithHTTPorHTTPS),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"settings": {
				Type:     schema.TypeList,
				Required: true,
				// Terraform isn't great at having Maps with complex subtypes. So making it a list with 1 item see also:
				// https://github.com/hashicorp/terraform-plugin-sdk/issues/62
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"icon": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"url": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: amplience.ValidateDiagWrapper(validation.IsURLWithHTTPorHTTPS),
									},
								},
							},
						},
						"visualization": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label": {
										Type:     schema.TypeString,
										Required: true,
									},
									"templated_uri": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: amplience.ValidateDiagWrapper(validation.IsURLWithHTTPorHTTPS),
									},
									"default": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceContentTypeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*amplience.Client)

	APIPath := fmt.Sprintf(c.ContentAPIPath+"/hubs/%[1]s/content-types", c.HubID)

	var contentType *amplience.RepositoryContentType
	var response *http.Response

	draft := createContentTypeDraft(data)

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
		return diag.FromErr(errorResponse)
	}

	if response == nil {
		return diag.FromErr(fmt.Errorf("could not create content type"))
	}

	err := amplience.ParseAndUnmarshalAmplienceResponseBody(response, &contentType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing response body into content type struct: %w", err))
	}
	data.SetId(contentType.ID)

	resourceContentTypeRead(ctx, data, meta)

	return diags
}

func createContentTypeDraft(d *schema.ResourceData) *amplience.RepositoryContentType {
	settings := resourceContentTypeGetSettings(d.Get("settings"))
	return &amplience.RepositoryContentType{
		ContentTypeURI: d.Get("content_type_uri").(string),
		Status:         d.Get("status").(string),
		Settings:       settings,
	}
}

func resourceContentTypeGetSettings(input interface{}) amplience.Settings {
	var result amplience.Settings

	// Below is a hacky workaround to allow for complex subtypes in a map like structure. We constrain the settings
	// field to be a list of max 1 item in the resource Schema and convert that to a map[string]interface{} below.
	inputSlice := input.([]interface{})
	i := inputSlice[0].(map[string]interface{})

	label, ok := i["label"].(string)
	if !ok {
		label = ""
	}
	var icons []amplience.Icon
	for _, rawIcon := range i["icon"].([]interface{}) {
		j := rawIcon.(map[string]interface{})
		size, ok := j["size"].(int)
		if !ok {
			size = 0
		}
		url, ok := j["url"].(string)
		if !ok {
			url = ""
		}
		icons = append(icons, amplience.Icon{
			Size: size,
			URL:  url,
		})
	}
	var visualizations []amplience.Visualization
	for _, rawVis := range i["visualization"].([]interface{}) {
		// TODO: if you get to the k's you know your method is too big. Cut this into chunks later
		// these type of functions happen so often and for almost all resources it might be worth it to create some
		// kind of resource field flattener function factory, this will make decent error handling easier as well
		k := rawVis.(map[string]interface{})
		vizLabel, ok := k["label"].(string)
		if !ok {
			vizLabel = ""
		}
		vizTempURI, ok := k["templated_uri"].(string)
		if !ok {
			vizTempURI = ""
		}
		vizDefault, ok := k["default"].(bool)
		if !ok {
			vizDefault = false
		}
		visualizations = append(visualizations, amplience.Visualization{
			Label:        vizLabel,
			TemplatedURI: vizTempURI,
			Default:      vizDefault,
		})
	}
	result = amplience.Settings{
		Label:          label,
		Icons:          icons,
		Visualizations: visualizations,
	}

	return result
}

func resourceContentTypeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	contentTypeID := data.Id()

	contentType, err := getContentTypeWithID(contentTypeID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred when trying to get content type with ID %s: %w", contentTypeID, err))
	}

	if contentType == nil {
		log.Print("[DEBUG] No content type found")
		data.SetId("")
	} else {
		log.Print("[DEBUG] Found following content type: ")
		log.Print(amplience.StringFormatObject(contentType))

		data.Set("content_type_uri", contentType.ContentTypeURI)
		data.Set("status", contentType.Status)
		settings := flattenContentTypeSettings(&contentType.Settings)
		data.Set("settings", settings)
	}
	return diags
}

func flattenContentTypeSettings(settings *amplience.Settings) interface{} {
	if settings != nil {
		st := make(map[string]interface{})
		st["label"] = settings.Label
		st["icons"] = flattenContentTypeSettingsIcons(&settings.Icons)
		st["visualizations"] = flattenContentTypeSettingsVisualizations(&settings.Visualizations)

		return st
	}
	return make(map[string]interface{})
}

func flattenContentTypeSettingsIcons(icons *[]amplience.Icon) []interface{} {
	if icons != nil {
		ics := make([]interface{}, len(*icons), len(*icons))

		for i, icon := range *icons {
			ic := make(map[string]interface{})

			ic["size"] = icon.Size
			ic["url"] = icon.URL

			ics[i] = ic
		}
		return ics
	}
	return make([]interface{}, 0)
}

func flattenContentTypeSettingsVisualizations(visualizations *[]amplience.Visualization) []interface{} {
	if visualizations != nil {
		vis := make([]interface{}, len(*visualizations), len(*visualizations))

		for i, visualization := range *visualizations {
			vi := make(map[string]interface{})

			vi["label"] = visualization.Label
			vi["default"] = visualization.Default
			vi["templated_uri"] = visualization.TemplatedURI

			vis[i] = vi
		}
		return vis
	}
	return make([]interface{}, 0)
}

func resourceContentTypeUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	contentTypeID := data.Id()

	if data.HasChange("settings") {
		draft := createContentTypeDraft(data)
		requestBody, err := json.Marshal(draft)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not marshal %v: %w", draft, err))
		}
		contentType, err := updateContentTypeWithID(contentTypeID, bytes.NewBuffer(requestBody), meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not update content type with ID %s: %w", contentTypeID, err))
		}
		if contentType == nil {
			log.Printf("[DEBUG] Nothing updated for content type with ID %s", contentTypeID)
		}
	}
	return diags
}

func resourceContentTypeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

// getContentTypeWithID sends a GET request to fetch a content type based on its ID
// The functionality of this should be abstracted into client library
func getContentTypeWithID(contentTypeID string, meta interface{}) (*amplience.RepositoryContentType, error) {
	contentType := amplience.RepositoryContentType{}

	c := meta.(*amplience.Client)
	APIPath := fmt.Sprintf(c.ContentAPIPath + "/content-types/" + contentTypeID)

	response, err := amplience.AmplienceRequest(APIPath, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make GET request to %s: %w", APIPath, err)
	}
	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &contentType)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into content type struct: %w", err)
	}
	return &contentType, nil
}

// updateContentTypeWithID sends a PATCH request to update a content type based on its ID
// The functionality of this should be abstracted into client library
func updateContentTypeWithID(contentTypeID string, requestBody *bytes.Buffer, meta interface{}) (*amplience.ContentType, error) {
	contentType := amplience.ContentType{}

	c := meta.(*amplience.Client)
	APIPath := fmt.Sprintf(c.ContentAPIPath + "/content-types/" + contentTypeID)
	response, err := amplience.AmplienceRequest(APIPath, http.MethodPatch, requestBody)
	if err != nil {
		return nil, fmt.Errorf("unable to make PATCH request to %s: %w", APIPath, err)
	}
	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &contentType)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into content type struct: %w", err)
	}
	return &contentType, nil
}
