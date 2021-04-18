package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentTypeCreate,
		ReadContext:   resourceContentTypeRead,
		UpdateContext: resourceContentTypeUpdate,
		DeleteContext: resourceContentTypeDelete,
		Schema: map[string]*schema.Schema{
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"content_type_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
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
	c := meta.(*content.Client)

	hub_id := data.Get("hub_id").(string)
	input := resourceContentTypeCreateInput(data)
	content_type, err := c.ContentTypeCreate(hub_id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, hub_id, content_type)
	return diags
}

func resourceContentTypeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	content_type_id := data.Id()
	hub_id := data.Get("hub_id").(string)
	content_type, err := c.ContentTypeGet(content_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, hub_id, content_type)
	return diags
}

func resourceContentTypeUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	content_type_id := data.Id()
	hub_id := data.Get("hub_id").(string)

	current, err := c.ContentTypeGet(content_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	input := resourceContentTypeCreateInput(data)
	content_type, err := c.ContentTypeUpdate(current, input)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeSaveState(data, hub_id, content_type)
	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentTypeDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*content.Client)

	id := data.Id()

	_, err := c.ContentTypeArchive(id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags
}

func resourceContentTypeSaveState(data *schema.ResourceData, hub_id string, resource content.ContentType) {
	data.SetId(resource.ID)
	data.Set("content_type_uri", resource.ContentTypeURI)
	data.Set("status", resource.Status)
	data.Set("label", resource.Settings.Label)
	data.Set("hub_id", hub_id)
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

func flattenContentTypeSettings(settings *content.ContentTypeSettings) interface{} {
	if settings != nil {
		st := make(map[string]interface{})
		st["label"] = settings.Label
		st["icons"] = flattenContentTypeSettingsIcons(&settings.Icons)
		st["visualizations"] = flattenContentTypeSettingsVisualizations(&settings.Visualizations)

		return st
	}
	return make(map[string]interface{})
}

func flattenContentTypeSettingsIcons(icons *[]content.ContentTypeIcon) []interface{} {
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

func flattenContentTypeSettingsVisualizations(visualizations *[]content.ContentTypeVisualization) []interface{} {
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