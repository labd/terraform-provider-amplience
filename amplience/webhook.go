package amplience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type WebhookEventsEnum string

const (
	WebhookContentItemAssigned        WebhookEventsEnum = "dynamic-content.content-item.assigned"
	WebhookContentItemCreated         WebhookEventsEnum = "dynamic-content.content-item.created"
	WebhookContentItemUpdated         WebhookEventsEnum = "dynamic-content.content-item.updated"
	WebhookContentItemWorkflowUpdated WebhookEventsEnum = "dynamic-content.content-item.workflow.updated"
	WebhookEditionPublished           WebhookEventsEnum = "dynamic-content.edition.published"
	WebhookEditionScheduled           WebhookEventsEnum = "dynamic-content.edition.scheduled"
	WebhookEditionUnscheduled         WebhookEventsEnum = "dynamic-content.edition.unscheduled"
	WebhookSnapshotPublished          WebhookEventsEnum = "dynamic-content.snapshot.published"
)

type Webhook struct {
	ID               string                `json:"id,omitempty"`
	Label            string                `json:"label"`
	Events           []string              `json:"events"`
	Handlers         []string              `json:"handlers"`
	Active           bool                  `json:"active"`
	Notifications    []Notification        `json:"notifications"`
	Secret           string                `json:"secret"`
	CreatedDate      *time.Time            `json:"createdDate,omitempty"`
	LastModifiedDate *time.Time            `json:"lastModifiedDate,omitempty"`
	Headers          []WebhookHeader       `json:"headers,omitempty"`
	Filters          []WebhookFilter       `json:"filters,omitempty"`
	Method           string                `json:"method"`
	CustomPayload    *WebhookCustomPayload `json:"customPayload,omitempty"`
}

type Notification struct {
	Email string
}

type WebhookHeader struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

type WebhookFilter struct {
	Type      string   `json:"type"`
	Arguments []RawArg `json:"arguments"`
}

type RawArg struct {
	JSONPath *string `json:"jsonPath,omitempty"`
	// We should never have a RawArg struct with both EqValue and InValues so we can set the json key to be equal
	InValues *[]string `json:"value,omitempty"`
	EqValue  *string   `json:"value,omitempty"`
}

type WebhookFilterInArguments struct {
	Value []string `json:"value"`
}

type WebhookFilterEqualArguments struct {
	Value string `json:"value"`
}

type WebhookFilterJSONPath struct {
	JSONPath string `json:"jsonPath"`
}

func (r *RawArg) UnmarshalJSON(buf []byte) error {
	jsonPath := WebhookFilterJSONPath{}
	err := json.Unmarshal(buf, &jsonPath)
	// If we find a jsonPath return immediately
	if err == nil && jsonPath.JSONPath != "" {
		r.JSONPath = &jsonPath.JSONPath
		return nil
	}
	if err != nil {
		return fmt.Errorf("error unmarshalling JSONPath: %w", err)
	}

	// Lord forgive me....
	// So if we get an array back as a response for the WebhookFilterArguments we need to Unmarshal the arguments into
	// a WebhookFilterInArguments and otherwise if it's a single string in an 'equal' filter we need a WebhookFilterEqualArguments
	// we check the bytes for a "[" to see if this is the case. Note that we know it's an argument and not a JSONPath
	// as otherwise we would have returned in the if clause above
	multiVal := bytes.Contains(buf, []byte("["))
	if multiVal {
		inArgs := WebhookFilterInArguments{}
		err = json.Unmarshal(buf, &inArgs)
		if err == nil && len(inArgs.Value) != 0 {
			r.InValues = &inArgs.Value
			return nil
		}
		if err != nil {
			return fmt.Errorf("error unmarshalling InValues: %w", err)
		}
	} else {
		eqArgs := WebhookFilterEqualArguments{}
		err = json.Unmarshal(buf, &eqArgs)
		if err == nil && eqArgs.Value != "" {
			r.EqValue = &eqArgs.Value
			return nil
		}
		if err != nil {
			return fmt.Errorf("error unmarshalling EqValue: %w", err)
		}
	}
	return nil
}

// The below custom Marshal func is a hacky solution to the problem that the API wants different types of values passed
// to it depending on the value of the "type" field in filter. This should be neatly abstracted away in an SDK
func (r RawArg) MarshalJSON() ([]byte, error) {
	if r.EqValue != nil && r.InValues != nil {
		return nil, fmt.Errorf("could not marshal %v, should not get both InValue and EqValue", r)
	}
	jsonPath, err := json.Marshal(r.JSONPath)
	if err == nil && len(jsonPath) > 0 && bytes.Compare(jsonPath, []byte("null")) != 0 {
		return []byte(fmt.Sprintf("{\"jsonPath\": %s }", jsonPath)), nil
	}
	if r.EqValue != nil {
		equalValue, err := json.Marshal(r.EqValue)
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf("{\"value\": %s }", equalValue)), nil
	} else if r.InValues != nil {
		inValues, err := json.Marshal(r.InValues)
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf("{\"value\": %s }", inValues)), nil
	}
	return []byte(""), nil
}

type WebhookCustomPayload struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
