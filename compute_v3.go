package conoha

import (
	"encoding/json"
	"fmt"
	"strings"

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
	if !res.IsStatus200() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return nil, fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	var v GetServersResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (api *V3) StartServer(serverId string) error {
	endpoint := api.Endpoints.Compute
	endpoint.Path = fmt.Sprintf(`/v2.1/servers/%s/action`, serverId)
	body := `{
		"os-start": null
	}`
	client := annette.New(endpoint)
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(strings.NewReader(body))
	if err != nil {
		return err
	}
	if !res.IsStatus202() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	return nil
}

func (api *V3) StopServer(serverId string) error {
	endpoint := api.Endpoints.Compute
	endpoint.Path = fmt.Sprintf(`/v2.1/servers/%s/action`, serverId)
	body := `{
		"os-stop": null
	}`
	client := annette.New(endpoint)
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(strings.NewReader(body))
	if err != nil {
		return err
	}
	if !res.IsStatus202() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	return nil
}

func (api *V3) RebootServer(serverId string) error {
	endpoint := api.Endpoints.Compute
	endpoint.Path = fmt.Sprintf(`/v2.1/servers/%s/action`, serverId)
	body := `{
		"reboot": {"type": "SOFT"}
	}`
	client := annette.New(endpoint)
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(strings.NewReader(body))
	if err != nil {
		return err
	}
	if !res.IsStatus202() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	return nil
}

func (api *V3) ForceShutdownServer(serverId string) error {
	endpoint := api.Endpoints.Compute
	endpoint.Path = fmt.Sprintf(`/v2.1/servers/%s/action`, serverId)
	body := `{
		"os-stop": {"force_shutdown": true}
	}`
	client := annette.New(endpoint)
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(strings.NewReader(body))
	if err != nil {
		return err
	}
	if !res.IsStatus202() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	return nil
}
