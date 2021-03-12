package amplience

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		Name   string
		Input  string
		Valid  []string
		Output bool
	}{
		{
			Name:   "Returns false if string not in slice",
			Input:  "an invalid string",
			Valid:  []string{"just valid strings here", "no invalid strings allowed"},
			Output: false,
		},
		{
			Name:   "Returns true if string in slice",
			Input:  "valid strings for life",
			Valid:  []string{"valid strings are the greatest strings", "valid strings for life"},
			Output: true,
		},
		{
			Name:   "Returns false if valid slice is empty",
			Input:  "",
			Valid:  []string{},
			Output: false,
		},
	}

	for _, tc := range tcs {
		tc := tc // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			got := StringInSlice(tc.Valid, tc.Input)
			if !assert.True(t, assert.ObjectsAreEqualValues(tc.Output, got)) {
				t.Logf("\n Got: %t \n Want: %t", got, tc.Output)
			}
		})
	}
}
