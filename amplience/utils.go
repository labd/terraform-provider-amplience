package amplience

import (
	"fmt"
	"strings"

	"github.com/labd/amplience-go-sdk/content"
)

func parseID(id string) (hub_id string, resource_id string) {
	values := strings.SplitN(id, ":", 2)
	if len(values) > 1 {
		return values[0], values[1]
	}
	return "", values[0]
}

func createID(hub_id string, resource_id string) string {
	return fmt.Sprintf("%s:%s", hub_id, resource_id)
}

type ClientInfo struct {
	Client *content.Client
	HubID  string
}

func getClient(meta interface{}) *ClientInfo {
	return meta.(*ClientInfo)
}
