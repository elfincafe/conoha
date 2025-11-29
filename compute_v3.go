package conoha

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	GetServerResponse struct {
		Server struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Status   string `json:"status"`
			TenantId string `json:"tenant_id"`
			UserId   string `json:"user_id"`
			Metadata struct {
				InstanceNameTag string `json:"instance_name_tag"`
				BackupStatus    string `json:"backup_status"`
				BackupId        string `json:"backup_id"`
				BackupSet       string `json:"backup_set"`
				BackupRotate    string `json:"backup_rotate"`
			} `json:"metadata"`
			HostId string `json:"hostId"`
			Image  string `json:"image"`
			Flavor struct {
				Id    string `json:"id"`
				Links []struct {
					Rel  string `json:"rel"`
					Href string `json:"href"`
				} `json:"links"`
			} `json:"flavor"`
			Created   time.Time `json:"created"`
			Updated   time.Time `json:"updated"`
			Addresses map[string]struct {
				Version         int    `json:"version"`
				Addr            string `json:"addr"`
				OsExtIpsType    string `json:"OS-EXT-IPS:type"`
				OsExtIpsMacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
			} `json:"addresses"`
			AccessIpv4 string `json:"accessIPv4"`
			AccessIpv6 string `json:"accessIPv6"`
			Links      []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
			OsDcfDiskConfig         string    `json:"OS-DCF:diskConfig"`
			OsExtAzAvailabilityZone string    `json:"OS-EXT-AZ:availability_zone"`
			ConfigDrive             string    `json:"config_drive"`
			KeyName                 string    `json:"key_name"`
			OsSrvUsgLaunchedAt      time.Time `json:"OS-SRV-USG:launched_at"`
			OsSrvUsgTeminatedAt     time.Time `json:"OS-SRV-USG:terminated_at"`
			SecurityGroups          []struct {
				Name string `json:"name"`
			} `json:"security_groups"`
			OsExtSrvAttrHost               string `json:"OS-EXT-SRV-ATTR:host"`
			OsExtSrvAttrInstanceName       string `json:"OS-EXT-SRV-ATTR:instance_name"`
			OsExtSrvAttrHypervisorHostname string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
			OsExtStsTaskState              string `json:"OS-EXT-STS:task_state"`
			OsExtStsVmState                string `json:"OS-EXT-STS:vm_state"`
			OsExtStsPowerState             int    `json:"OS-EXT-STS:power_state"`
			OsExtendedVolumesAttached      []struct {
				Id string `json:"id"`
			} `json:"os-extended-volumes:volumes_attached"`
		} `json:"server"`
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

func (api *V3) GetServer(id string) (*GetServerResponse, error) {
	endpoint := api.Endpoints.Compute
	endpoint.Path = fmt.Sprintf(`/v2.1/servers/%s`, id)
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
	var v GetServerResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
