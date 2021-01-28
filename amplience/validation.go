package amplience

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValidateDiagWrapper wraps a deprecated schema.ValidateFunc and returns a schema.SchemaValidateDiagFunc
// this is needed as a placeholder until the base validation funcs in helper/validation are updated
// see: https://discuss.hashicorp.com/t/validatefunc-deprecation-in-terraform-plugin-sdk-v2/12000
func ValidateDiagWrapper(validateFunc func(interface{}, string) ([]string, []error)) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		warnings, errs := validateFunc(i, fmt.Sprintf("%+v", path))
		var diags diag.Diagnostics
		for _, warning := range warnings {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  warning,
			})
		}
		for _, err := range errs {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		return diags
	}
}
