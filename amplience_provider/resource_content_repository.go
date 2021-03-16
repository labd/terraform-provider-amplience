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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/labd/terraform-provider-amplience/amplience"
)

func resourceContentRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentRepositoryCreate,
		ReadContext:   resourceContentRepositoryRead,
		UpdateContext: resourceContentRepositoryUpdate,
		DeleteContext: resourceContentRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringDoesNotContainAny(" "),
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_types": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hub_content_type_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content_type_uri": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"item_locales": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceContentRepositoryCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*amplience.ClientConfig)

	apiPath := fmt.Sprintf(c.ContentApiUrl+"/hubs/%[1]s/content-repositories", c.HubID)

	var repository *amplience.ContentRepository
	var response *http.Response

	draft := createContentRepositoryDraft(data)

	errorResponse := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		var err error
		requestBody, err := json.Marshal(draft)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error marshalling draft %v: %w", draft, err))
		}
		response, err = amplience.AmplienceRequest(c, apiPath, http.MethodPost, bytes.NewBuffer(requestBody))
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error during http request: %w", err))
		}
		return amplience.HandleAmplienceError(response)
	})

	if errorResponse != nil {
		return diag.FromErr(errorResponse)
	}

	if response == nil {
		return diag.FromErr(fmt.Errorf("could not create content repository"))
	}
	err := amplience.ParseAndUnmarshalAmplienceResponseBody(response, &repository)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing response body into repository struct: %w", err))
	}
	data.SetId(repository.ID)

	resourceContentRepositoryRead(ctx, data, meta)

	return diags
}

func createContentRepositoryDraft(data *schema.ResourceData) *amplience.ContentRepository {
	contentTypes := resourceContentRepositoryGetContentTypes(data.Get("content_types"))

	var featureSlice []string
	for _, val := range data.Get("features").([]interface{}) {
		feature := val.(string)
		featureSlice = append(featureSlice, feature)
	}

	var itemLocaleSlice []string
	for _, val := range data.Get("item_locales").([]interface{}) {
		itemLocale := val.(string)
		itemLocaleSlice = append(itemLocaleSlice, itemLocale)
	}

	return &amplience.ContentRepository{
		Name:         data.Get("name").(string),
		Label:        data.Get("label").(string),
		Status:       data.Get("status").(string),
		Features:     featureSlice,
		Type:         data.Get("type").(string),
		ContentTypes: contentTypes,
		ItemLocales:  itemLocaleSlice,
	}
}

func resourceContentRepositoryRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	repositoryID := data.Id()

	repository, err := getContentRepositoryWithID(repositoryID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred when trying to get content repository with ID %s: %w", repositoryID, err))
	}

	if repository == nil {
		log.Print("[DEBUG] No content repository found")
		data.SetId("")
	} else {
		log.Print("[DEBUG] Found following content repository: ")
		log.Print(amplience.StringFormatObject(repository))

		data.Set("name", repository.Name)
		data.Set("label", repository.Label)
		data.Set("status", repository.Status)
		data.Set("features", repository.Features)
		data.Set("type", repository.Type)
		contentTypes := flattenContentRepositoryContentTypes(&repository.ContentTypes)
		data.Set("content_types", contentTypes)
		data.Set("item_locales", repository.ItemLocales)
	}
	return diags
}

func flattenContentRepositoryContentTypes(contentTypes *[]amplience.ContentType) []interface{} {
	if contentTypes != nil {
		cts := make([]interface{}, len(*contentTypes), len(*contentTypes))

		for i, contentType := range *contentTypes {
			ct := make(map[string]interface{})

			ct["hub_content_type_id"] = contentType.HubContentTypeID
			ct["content_type_uri"] = contentType.ContentTypeURI
			cts[i] = ct
		}
		return cts
	}
	return make([]interface{}, 0)
}

func resourceContentRepositoryUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	repositoryID := data.Id()

	if data.HasChange("label") {
		draft := createContentRepositoryDraft(data)
		requestBody, err := json.Marshal(draft)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not marshal %v: %w", draft, err))
		}

		repository, err := updateContentRepositoryWithID(repositoryID, bytes.NewBuffer(requestBody), meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not update repository with ID %s: %w", repositoryID, err))
		}
		if repository == nil {
			log.Printf("[DEBUG] Nothing updated for repository with ID %s", repositoryID)
		}
	}

	return resourceContentRepositoryRead(ctx, data, meta)
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentRepositoryDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	data.SetId("")
	return diags
}

// getContentRepositoryWithID returns a content repository based on its ID if it exists. Else returns nil.
// The functionality of this should be abstracted into client library
func getContentRepositoryWithID(contentRepoID string, meta interface{}) (*amplience.ContentRepository, error) {
	repository := amplience.ContentRepository{}

	c := meta.(*amplience.ClientConfig)
	apiPath := fmt.Sprintf(c.ContentApiUrl + "/content-repositories/" + contentRepoID)

	response, err := amplience.AmplienceRequest(c, apiPath, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make GET request to %s: %w", apiPath, err)
	}

	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &repository)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into repository struct: %w", err)
	}
	return &repository, nil
}

// updateContentRepositoryWithID sends a PATCH request to update a content repository based on its ID
// The functionality of this should be abstracted into client library
func updateContentRepositoryWithID(contentRepoID string, requestBody *bytes.Buffer, meta interface{}) (*amplience.ContentRepository, error) {
	repository := amplience.ContentRepository{}

	c := meta.(*amplience.ClientConfig)
	apiPath := fmt.Sprintf(c.ContentApiUrl + "/content-repositories/" + contentRepoID)

	response, err := amplience.AmplienceRequest(c, apiPath, http.MethodPatch, requestBody)
	if err != nil {
		return nil, fmt.Errorf("unable to make GET request to %s: %w", apiPath, err)
	}
	err = amplience.ParseAndUnmarshalAmplienceResponseBody(response, &repository)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into repository struct: %w", err)
	}
	return &repository, nil
}

func resourceContentRepositoryGetContentTypes(input interface{}) []amplience.ContentType {
	inputSlice := input.([]interface{})
	var result []amplience.ContentType

	for _, raw := range inputSlice {
		i := raw.(map[string]interface{})

		hubContentTypeID, ok := i["hub_content_type_id"].(string)
		if !ok {
			hubContentTypeID = ""
		}
		contentTypeURI, ok := i["content_type_uri"].(string)
		if !ok {
			contentTypeURI = ""
		}
		result = append(result, amplience.ContentType{
			HubContentTypeID: hubContentTypeID,
			ContentTypeURI:   contentTypeURI,
		})
	}

	return result
}
