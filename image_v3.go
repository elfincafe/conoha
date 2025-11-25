package conoha

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elfincafe/annette"
	"github.com/google/uuid"
)

type (
	CreateIsoImageResponse struct {
		HwRescueBus     string    `json:"hw_rescue_bus"`
		HwRescueDevice  string    `json:"hw_rescue_device"`
		Name            string    `json:"name"`
		DiskFormat      string    `json:"disk_format"`
		ContainerFormat string    `json:"container_format"`
		Visibility      string    `json:"visibility"`
		Size            int       `json:"size"`
		VirtualSize     int       `json:"virtual_size"`
		Status          string    `json:"status"`
		Checksum        int       `json:"checksum"`
		Protected       bool      `json:"protected"`
		MinRam          int       `json:"min_ram"`
		MinDisk         int       `json:"min_disk"`
		Owner           string    `json:"owner"`
		OsHidden        bool      `json:"os_hidden"`
		OsHashAlgo      string    `json:"os_hash_algo"`
		OsHashValue     string    `json:"os_hash_value"`
		Id              string    `json:"id"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		Tags            []string  `json:"tags"`
		Self            string    `json:"self"`
		File            string    `json:"file"`
		Schema          string    `json:"schema"`
	}
)

func (api *V3) CreateIsoImage(name string) (*CreateIsoImageResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = "/v2/images"
	body := `
		{
			"name": "###NAME###",
			"disk_format": "iso",
			"hw_rescue_bus": "ide",
			"hw_rescue_device": "cdrom",
			"container_format": "bare"
		}
  	`
	if len(name) == 0 {
		u, _ := uuid.NewRandom()
		name = u.String()
	}
	name = ""
	body = strings.ReplaceAll(body, "###NAME###", name)
	client := annette.New(endpoint)
	client.SetHeader("X-Auth-Token", api.Token)
	res, err := client.Post(body)
	if err != nil {
		return nil, err
	}
	if !res.IsStatus201() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return nil, fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	var v CreateIsoImageResponse
	content := res.Binary()
	err = json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, err
}
