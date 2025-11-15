package k8s

import (
	"encoding/json"

	"k8s.io/client-go/discovery"
)

func ListServerGroupsAndResources(dc discovery.DiscoveryInterface) (string, error) {
	groups, resources, err := dc.ServerGroupsAndResources()
	if err != nil {
		return "", err
	}

	data := struct {
		Groups    interface{} `json:"groups"`
		Resources interface{} `json:"resources"`
	}{
		Groups:    groups,
		Resources: resources,
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
