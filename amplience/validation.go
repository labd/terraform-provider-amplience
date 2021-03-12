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

// StringInSlice takes a slice and looks for an element in it. If found it will return true
// can be used to manually validate elements in a list in the create/update process as terraform validation functions
// are not designed for lists
func StringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
