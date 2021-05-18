package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentTypeAssignment() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource used to assign a Content Type to a Content Repository",
		CreateContext: resourceContentTypeAssignmentCreate,
		ReadContext:   resourceContentTypeAssignmentRead,
		DeleteContext: resourceContentTypeAssignmentDelete,
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Description: "ID of the Content Repository to assign the type to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"content_type_id": {
				Description: "ID of the Content Type to assign to the Repository",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceContentTypeAssignmentCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	repository_id := data.Get("repository_id").(string)
	content_type_id := data.Get("content_type_id").(string)

	_, err := ci.client.ContentRepositoryAssignContentType(repository_id, content_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceContentTypeAssignmentSaveState(data, repository_id, content_type_id)
	return diags
}

func resourceContentTypeAssignmentRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	repository_id, content_type_id := parseID(data.Id())
	// TODO: check amplience
	resourceContentTypeAssignmentSaveState(data, repository_id, content_type_id)

	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentTypeAssignmentDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	ci := getClient(meta)

	repository_id, content_type_id := parseID(data.Id())

	_, err := ci.client.ContentRepositoryRemoveContentType(repository_id, content_type_id)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return diags

}

func resourceContentTypeAssignmentSaveState(data *schema.ResourceData, repository_id string, content_type_id string) {
	data.SetId(createID(repository_id, content_type_id))
	data.Set("repository_id", repository_id)
	data.Set("content_type_id", content_type_id)
}
