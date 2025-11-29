package conoha

import (
	"encoding/json"

	"github.com/elfincafe/annette"
)

type (
	GetServersResponse struct {
		Servers []struct {
			Id    string `json:"id"`
			Name  string `json:"name"`
			Links []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"servers"`
	}
)

func (api *V3) GetServers() (*GetServersResponse, error) {
	endpoint := api.Endpoints.Compute
	endpoint.Path = "/v2.1/servers"
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Get()
	if err != nil {
		return nil, err
	}
	var v GetServersResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
