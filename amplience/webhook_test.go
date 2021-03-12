package amplience

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawArgUnmarshalJson(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		Name             string
		Input            []byte
		ExpectedError    bool
		ExpectedJSONPath *string
		ExpectedInValues *[]string
		ExpectedEqValue  *string
	}{
		{
			Name:             "Unmarshals JSONPath if only JSON path present",
			Input:            []byte("{\n      \"jsonPath\" : \"some path\"\n    }"),
			ExpectedError:    false,
			ExpectedJSONPath: createPointer("some path"),
			ExpectedInValues: nil,
			ExpectedEqValue:  nil,
		},
		{
			Name:             "Unmarshals only JSONPath if JSON path and other stuff is present",
			Input:            []byte("{\n    \"jsonPath\": \"some path\",\n    \"value\": [\"value1\", \"value2\"],\n    \"random_other_field\": 12345\n}"),
			ExpectedError:    false,
			ExpectedJSONPath: createPointer("some path"),
			ExpectedInValues: nil,
			ExpectedEqValue:  nil,
		},
		{
			Name:             "Unmarshals only InValues if only values: []string is present ",
			Input:            []byte("{\n      \"value\" : [\"an_in_value\", \"some_other_in_value\"] \n}"),
			ExpectedError:    false,
			ExpectedJSONPath: nil,
			ExpectedInValues: &[]string{"an_in_value", "some_other_in_value"},
			ExpectedEqValue:  nil,
		},
		{
			Name:             "Unmarshals only InValues if values: []string and other stuff is present ",
			Input:            []byte("{\n  \"value\": [\"an_in_value\", \"some_other_in_value\"],\n    \"random_other_field\": \"12345\"\n}"),
			ExpectedError:    false,
			ExpectedJSONPath: nil,
			ExpectedInValues: &[]string{"an_in_value", "some_other_in_value"},
			ExpectedEqValue:  nil,
		},
		{
			Name:             "Unmarshals only EqValue if only value: string is present",
			Input:            []byte("{\n      \"value\" : \"an_eq_value\" \n}"),
			ExpectedError:    false,
			ExpectedJSONPath: nil,
			ExpectedInValues: nil,
			ExpectedEqValue:  createPointer("an_eq_value"),
		},
		{
			Name:             "Unmarshals only EqValue if value: string and other stuff is present ",
			Input:            []byte("{\n  \"value\": \"an_eq_value\",\n    \"random_other_field\": \"12345\"\n}"),
			ExpectedError:    false,
			ExpectedJSONPath: nil,
			ExpectedInValues: nil,
			ExpectedEqValue:  createPointer("an_eq_value"),
		},
		{
			Name:             "Returns error if both value: []string and value: string is passed",
			Input:            []byte("{\n  \"value\": [\"an_in_value\", \"some_other_in_value\"],\n \"value\": \"an_eq_value\",\n    \"random_other_field\": \"12345\"\n}"),
			ExpectedError:    true,
			ExpectedJSONPath: nil,
			ExpectedInValues: nil,
			ExpectedEqValue:  nil,
		},
		{
			Name:             "Returns error if invalid json",
			Input:            []byte("{\n  \"value\": [\"an_in_value\", \"some_other_in_value\"],\n \"value\": \"an_eq_value\",\n    \"random_other_field\": \"12345\",\n}"),
			ExpectedError:    true,
			ExpectedJSONPath: nil,
			ExpectedInValues: nil,
			ExpectedEqValue:  nil,
		},
	}

	for _, tc := range tcs {
		tc := tc // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			arg := RawArg{}
			err := arg.UnmarshalJSON(tc.Input)
			if tc.ExpectedError {
				assert.True(t, assert.NotNil(t, err))
			}
			assert.True(t, assert.ObjectsAreEqualValues(tc.ExpectedJSONPath, arg.JSONPath))
			assert.True(t, assert.ObjectsAreEqualValues(tc.ExpectedInValues, arg.InValues))
			assert.True(t, assert.ObjectsAreEqualValues(tc.ExpectedEqValue, arg.EqValue))

		})
	}
}

func TestRawArgMarshalJSON(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		Name           string
		Input          RawArg
		ExpectedError  bool
		ExpectedResult []byte
	}{
		{
			Name: "Correctly marshals jsonPath if only jsonPath present",
			Input: RawArg{
				JSONPath: createPointer("bla_path"),
				InValues: nil,
				EqValue:  nil,
			},
			ExpectedError:  false,
			ExpectedResult: []byte(fmt.Sprintf("{\"jsonPath\": \"bla_path\" }")),
		},
		{
			Name: "Correctly marshals jsonPath if jsonPath and other stuff present",
			Input: RawArg{
				JSONPath: createPointer("bla_path"),
				InValues: &[]string{"other_things", "in_value"},
				EqValue:  nil,
			},
			ExpectedError:  false,
			ExpectedResult: []byte(fmt.Sprint("{\"jsonPath\": \"bla_path\" }")),
		},
		{
			Name: "Correctly marshals InValues if only InValues present",
			Input: RawArg{
				JSONPath: nil,
				InValues: &[]string{"other_things", "in_value"},
				EqValue:  nil,
			},
			ExpectedError:  false,
			ExpectedResult: []byte(fmt.Sprint("{\"value\": [\"other_things\",\"in_value\"] }")),
		},
		{
			Name: "Correctly marshals EqValue if only EqValue present",
			Input: RawArg{
				JSONPath: nil,
				InValues: nil,
				EqValue:  createPointer("equal_value"),
			},
			ExpectedError:  false,
			ExpectedResult: []byte(fmt.Sprint("{\"value\": \"equal_value\" }")),
		},
		{
			Name: "Marshals into empty byte slice if no values set",
			Input: RawArg{
				JSONPath: nil,
				InValues: nil,
				EqValue:  nil,
			},
			ExpectedError:  false,
			ExpectedResult: []byte(fmt.Sprint("")),
		},
		{
			Name: "Returns error if both EqValue and InValues fields are set",
			Input: RawArg{
				JSONPath: createPointer("bla_path"),
				InValues: &[]string{"in_value", "is_set"},
				EqValue:  createPointer("is_also_set"),
			},
			ExpectedError:  true,
			ExpectedResult: nil,
		},
	}
	for _, tc := range tcs {
		tc := tc // https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			bytes, err := tc.Input.MarshalJSON()
			if tc.ExpectedError {
				assert.True(t, assert.NotNil(t, err))
			}
			if !assert.True(t, assert.ObjectsAreEqualValues(tc.ExpectedResult, bytes)) {
				t.Logf("\n Got: %s \n Want: %s", bytes, tc.ExpectedResult)
			}

		})
	}
}

func createPointer(s string) *string {
	result := fmt.Sprint(s)
	return &result
}
