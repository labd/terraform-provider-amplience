package amplience

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHandleAmplienceError(t *testing.T) {
	t.Parallel()

	mockErrorResponse := ErrorResponse{Errors: []ErrorResponseItem{
		{
			Level:    "Error",
			Code:     "SOME_CODE",
			Message:  "AN ERROR MESSAGE",
			Property: "Props",
		},
		{
			Level:        "Error",
			Code:         "SOME_OTHER_CODE",
			Message:      "AN ERROR MESSAGE",
			InvalidValue: "42",
			Entity:       "ENTITY SCHMENTITY",
		},
	}}
	tcs := []struct {
		Name   string
		Input  *http.Response
		Output *resource.RetryError
	}{
		{
			Name: "Returns nil if StatusCode is 200",
			Input: &http.Response{
				StatusCode: 200,
			},
			Output: nil,
		},
		{
			Name:   "Returns error if Input is nil",
			Input:  nil,
			Output: resource.NonRetryableError(fmt.Errorf("received nil response, unable to handle error")),
		},
		{
			Name: "Returns retryable error message on StatusCode 500",
			Input: &http.Response{
				StatusCode: 500,
				Status:     "RetryError",
				Body:       ioutil.NopCloser(strings.NewReader(StringFormatObject(mockErrorResponse))),
			},
			Output: resource.RetryableError(fmt.Errorf("retryable error with code 500 received: RetryError\n Amplience Error Response: %s", StringFormatObject(mockErrorResponse))),
		},
		{
			Name: "Returns a non-retryable error message on StatusCode 400",
			Input: &http.Response{
				StatusCode: 400,
				Status:     "NonRetryError",
				Body:       ioutil.NopCloser(strings.NewReader(StringFormatObject(mockErrorResponse))),
			},
			Output: resource.NonRetryableError(fmt.Errorf("non retryable error with code 400 received: NonRetryError\n Amplience Error Response: %s", StringFormatObject(mockErrorResponse))),
		},
	}

	for _, tc := range tcs {
		tc := tc // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			got := HandleAmplienceError(tc.Input)
			if !assert.True(t, assert.ObjectsAreEqualValues(tc.Output, got)) {
				t.Logf("\n Got: %t \n Want: %t", got, tc.Output)
			}
		})
	}
}
