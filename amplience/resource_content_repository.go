package amplience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/labd/amplience-go-sdk/content"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentRepositoryCreate,
		ReadContext:   resourceContentRepositoryRead,
		UpdateContext: resourceContentRepositoryUpdate,
		DeleteContext: resourceContentRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateDiagWrapper(validation.StringDoesNotContainAny(" ")),
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceContentRepositoryCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	hub_id := data.Get("hub_id").(string)
	input := content.ContentRepositoryInput{
		Name:  data.Get("name").(string),
		Label: data.Get("label").(string),
	}

	repository, err := c.ContentRepositoryCreate(hub_id, input)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(repository.ID)
	data.Set("name", repository.Name)
	data.Set("label", repository.Label)
	data.Set("hub_id", hub_id)
	return diags
}

func resourceContentRepositoryRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	repository_id := data.Id()

	repository, err := c.ContentRepositoryGet(repository_id)
	if err != nil {
		return diag.FromErr(err)
	}

	hub, err := repository.GetHub(c)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("name", repository.Name)
	data.Set("label", repository.Label)
	data.Set("hub_id", hub.ID)
	return diags
}

func resourceContentRepositoryUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*content.Client)

	repository_id := data.Id()

	if data.HasChange("label") || data.HasChange("name") {
		current, err := c.ContentRepositoryGet(repository_id)
		if err != nil {
			return diag.FromErr(err)
		}

		input := content.ContentRepositoryInput{
			Name:  data.Get("name").(string),
			Label: data.Get("label").(string),
		}

		repository, err := c.ContentRepositoryUpdate(current, input)
		if err != nil {
			return diag.FromErr(err)
		}

		data.Set("label", repository.Label)
		data.Set("name", repository.Name)
	}

	return diags
}

// The amplience API does not have a repository delete functionality. Setting ID to "" and returning nil
func resourceContentRepositoryDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	data.SetId("")
	return diags
}
