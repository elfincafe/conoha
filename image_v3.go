package conoha

import (
	"encoding/json"
	"fmt"
	"os"
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
	body := fmt.Sprintf(`{
		"name": "%s",
		"disk_format": "iso",
		"hw_rescue_bus": "ide",
		"hw_rescue_device": "cdrom",
		"container_format": "bare"
  	}`, name)
	if name == "" {
		u, _ := uuid.NewRandom()
		name = u.String()
	}
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !res.IsStatus201() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return nil, fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	var v CreateIsoImageResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, err
}

func (api *V3) UploadIsoImage(imageId, path string) error {
	endpoint := api.Endpoints.Image
	endpoint.Path = fmt.Sprintf("/v2/images/%s/file", imageId)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("Content-Type", "application/octet-stream")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.UploadByPut(f)
	if err != nil {
		return err
	}
	if !res.IsStatus204() {
		var v ConohaError
		json.Unmarshal(res.Binary(), &v)
		return fmt.Errorf("status:%d, error:%s", v.Code, v.Error)
	}
	return nil
}
