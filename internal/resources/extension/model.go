package extension

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/amplience-go-sdk/content"
)

type Extension struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Label                     types.String `tfsdk:"label"`
	Description               types.String `tfsdk:"description"`
	URL                       types.String `tfsdk:"url"`
	Height                    types.Int64  `tfsdk:"height"`
	EnabledForAllContentTypes types.Bool   `tfsdk:"enabled_for_all_content_types"`
	Category                  types.String `tfsdk:"category"`
	Parameters                types.String `tfsdk:"parameters"`
	Snippets                  []Snippet    `tfsdk:"snippets"`
	Settings                  types.String `tfsdk:"settings"`
}

type Snippet struct {
	Label types.String `tfsdk:"label"`
	Body  types.String `tfsdk:"body"`
}

func (e *Extension) ToInput(diags *diag.Diagnostics) content.ExtensionInput {
	input := content.ExtensionInput{
		Name:                      e.Name.ValueString(),
		Label:                     e.Label.ValueString(),
		Description:               e.Description.ValueString(),
		URL:                       e.URL.ValueString(),
		Height:                    int(e.Height.ValueInt64()),
		EnabledForAllContentTypes: e.EnabledForAllContentTypes.ValueBool(),
		Category:                  e.Category.ValueString(),
		Settings:                  e.Settings.ValueString(),
	}

	input.Parameters = e.Parameters.ValueString()

	for _, s := range e.Snippets {
		input.Snippets = append(input.Snippets, content.ExtensionSnippet{
			Label: s.Label.ValueString(),
			Body:  s.Body.ValueString(),
		})
	}

	return input
}

func NewExtensionFromNative(e *content.Extension, diags *diag.Diagnostics) *Extension {
	ext := &Extension{
		ID:                        types.StringValue(e.ID),
		Name:                      types.StringValue(e.Name),
		Label:                     types.StringValue(e.Label),
		Description:               nullableString(e.Description),
		URL:                       types.StringValue(e.URL),
		Height:                    types.Int64Value(int64(e.Height)),
		EnabledForAllContentTypes: types.BoolValue(e.EnabledForAllContentTypes),
		Category:                  types.StringValue(e.Category),
		Settings:                  nullableString(e.Settings),
		Parameters:                nullableString(e.Parameters),
		Snippets:                  nil,
	}

	for _, s := range e.Snippets {
		ext.Snippets = append(ext.Snippets, Snippet{
			Label: types.StringValue(s.Label),
			Body:  types.StringValue(s.Body),
		})
	}

	return ext
}

func nullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func nullableInt64(i int) types.Int64 {
	if i == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(int64(i))
}
