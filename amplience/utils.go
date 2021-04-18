package amplience

import (
	"fmt"
	"strings"
)

func parseID(id string) (hub_id string, resource_id string) {
	values := strings.SplitN(id, ":", 2)

	return values[0], values[1]
}

func createID(hub_id string, resource_id string) string {
	return fmt.Sprintf("%s:%s", hub_id, resource_id)
}
