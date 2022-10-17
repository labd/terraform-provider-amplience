package amplience

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentType() *schema.Resource {
	return &schema.Resource{
		Description: "Content types are the templates for content items, defining the type of content to be created, " +
			"including its structure and validation rules. Content types are stored externally to Dynamic Content, " +
			"on web based services such as AWS, and must be registered with a hub before they can be used to create " +
			"content.\n" +
			"For more info see [Amplience Content Type Docs](https://amplience.com/docs/integration/workingwithcontenttypes.html)",
		CreateContext: resourceContentTypeCreate,
		ReadContext:   resourceContentTypeRead,
		UpdateContext: resourceContentTypeUpdate,
		DeleteContext: resourceContentTypeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"content_type_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Description: "Status of the Content Type. Can be ACTIVE or DELETED",
				Type:        schema.TypeString,
				Required:    true,
			},
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
							ValidateDiagFunc: ValidateDiagWrapper(validation.IsURLWithHTTPorHTTPS),
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
							ValidateDiagFunc: ValidateDiagWrapper(validation.IsURLWithHTTPorHTTPS),
						},
						"default": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceContentTypeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	input := resourceContentTypeCreateInput(data)
	instance, err := ci.client.ContentTypeCreate(ci.hubID, input)

	if errResp, ok := err.(*content.ErrorResponse); ok {
		if errResp.StatusCode >= 400 {

			log.Println("Received 400 conflict response: content type already exists.")
			log.Println("Proceeding to unarchive if necessary and update exiting content type.")

			instance, err = ci.client.ContentTypeFindByUri(input.ContentTypeURI, ci.hubID)

			if err != nil {
				return diag.FromErr(err)
			}

			if instance.Status == string(content.StatusArchived) {
				instance, err = ci.client.ContentTypeUnarchive(instance.ID)
			}

			instance, err = ci.client.ContentTypeUpdate(instance, input)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, instance)
	return diags
}

func resourceContentTypeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	content_type_id := data.Id()
	content_type, err := ci.client.ContentTypeGet(content_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, content_type)
	return diags
}

func resourceContentTypeUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()

	instance, err := ci.client.ContentTypeGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if instance.Status == string(content.StatusArchived) {
		log.Println("Content type was archived. Proceed to unarchive first before applying update.")

		instance, err = ci.client.ContentTypeUnarchive(instance.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	input := resourceContentTypeCreateInput(data)
	content_type, err := ci.client.ContentTypeUpdate(instance, input)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, content_type)
	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentTypeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	id := data.Id()

	_, err := ci.client.ContentTypeArchive(id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceContentTypeSaveState(data *schema.ResourceData, resource content.ContentType) {
	icons := marshallContentTypeSettingsIcons(&resource.Settings.Icons)
	visualizations := marshallContentTypeSettingsVisualizations(&resource.Settings.Visualizations)

	data.SetId(resource.ID)
	data.Set("content_type_uri", resource.ContentTypeURI)
	data.Set("status", resource.Status)
	data.Set("label", resource.Settings.Label)
	data.Set("visualization", visualizations)
	data.Set("icon", icons)
}

func resourceContentTypeCreateInput(data *schema.ResourceData) content.ContentTypeInput {
	settings := resourceContentTypeGetSettings(data)

	return content.ContentTypeInput{
		ContentTypeURI: data.Get("content_type_uri").(string),
		Settings:       settings,
	}
}

func resourceContentTypeGetSettings(data *schema.ResourceData) content.ContentTypeSettings {
	var result content.ContentTypeSettings

	label := data.Get("label").(string)
	var icons []content.ContentTypeIcon
	for _, rawIcon := range data.Get("icon").([]interface{}) {
		j := rawIcon.(map[string]interface{})
		size, ok := j["size"].(int)
		if !ok {
			size = 0
		}
		url, ok := j["url"].(string)
		if !ok {
			url = ""
		}
		icons = append(icons, content.ContentTypeIcon{
			Size: size,
			URL:  url,
		})
	}
	var visualizations []content.ContentTypeVisualization
	for _, rawVis := range data.Get("visualization").([]interface{}) {
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
		visualizations = append(visualizations, content.ContentTypeVisualization{
			Label:        vizLabel,
			TemplatedURI: vizTempURI,
			Default:      vizDefault,
		})
	}
	result = content.ContentTypeSettings{
		Label:          label,
		Icons:          icons,
		Visualizations: visualizations,
	}

	return result
}

func marshallContentTypeSettingsIcons(icons *[]content.ContentTypeIcon) []interface{} {
	if icons != nil {
		ics := make([]interface{}, len(*icons))

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

func marshallContentTypeSettingsVisualizations(visualizations *[]content.ContentTypeVisualization) []interface{} {
	if visualizations != nil {
		vis := make([]interface{}, len(*visualizations))

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
