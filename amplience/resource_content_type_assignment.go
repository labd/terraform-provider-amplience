package amplience

import (
	"context"

	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentTypeAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentTypeAssignmentCreate,
		ReadContext:   resourceContentTypeAssignmentRead,
		DeleteContext: resourceContentTypeAssignmentDelete,
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceContentTypeAssignmentCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	repository_id := data.Get("repository_id").(string)
	content_type_id := data.Get("content_type_id").(string)

	_, err := c.ContentRepositoryAssignContentType(repository_id, content_type_id)
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
func resourceContentTypeAssignmentDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	data.SetId("")
	return diags
}

func resourceContentTypeAssignmentSaveState(data *schema.ResourceData, repository_id string, content_type_id string) {
	data.SetId(createID(repository_id, content_type_id))
	data.Set("repository_id", repository_id)
	data.Set("content_type_id", content_type_id)
}
