package amplience

import (
	"encoding/json"
	"fmt"
)

func StringFormatObject(object interface{}) string {
	data, err := json.MarshalIndent(object, "", "    ")

	if err != nil {
		return fmt.Sprintf("%+v", object)
	}
	return string(append(data, '\n'))
}
