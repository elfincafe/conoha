package conoha

import (
	"encoding/json"
	"fmt"
	"net/url"
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
		Id              uuid.UUID `json:"id"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
		Tags            []string  `json:"tags"`
		Self            string    `json:"self"`
		File            string    `json:"file"`
		Schema          string    `json:"schema"`
	}
	GetImagesResponse struct {
		Images []struct {
			Status          string    `json:"status"`
			Name            string    `json:"name"`
			Tags            []string  `json:"tags"`
			ContainerFormat string    `json:"container_format"`
			CreatedAt       time.Time `json:"created_at"`
			DiskFormat      string    `json:"disk_format"`
			UpdatedAt       time.Time `json:"updated_at"`
			Visibility      string    `json:"visibility"`
			Self            string    `json:"self"`
			MinDisk         int       `json:"min_disk"`
			Protected       bool      `json:"protected"`
			Id              uuid.UUID `json:"id"`
			File            string    `json:"file"`
			Checksum        string    `json:"checksum"`
			OsType          string    `json:"os_type"`
			OsHashAlgo      string    `json:"os_hash_algo"`
			OsHashValue     string    `json:"os_hash_value"`
			OsHidden        bool      `json:"os_hidden"`
			Owner           string    `json:"owner"`
			Size            int       `json:"size"`
			MinRam          int       `json:"min_ram"`
			Schema          string    `json:"schema"`
			VirtualSize     int       `json:"virtual_size"`
			Architecture    string    `json:"architecture"`
		} `json:"images"`
		Schema string `json:"schema"`
		First  string `json:"first"`
	}
)

func (api *V3) UploadIsoImage(imageId uuid.UUID, path string) error {
	endpoint := api.Endpoints.Image
	endpoint.Path = fmt.Sprintf("/v2/images/%s/file", imageId.String())
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

func (api *V3) CreateIsoImage(name string) (*CreateIsoImageResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = "/v2/images"
	if name == "" {
		u, _ := uuid.NewRandom()
		name = u.String()
	}
	body := fmt.Sprintf(`{
		"name": "%s",
		"disk_format": "iso",
		"hw_rescue_bus": "ide",
		"hw_rescue_device": "cdrom",
		"container_format": "bare"
  	}`, name)
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

func (api *V3) GetImages(args map[string]string) (*GetImagesResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = "/v2/images"
	q := url.Values{}
	for k, v := range args {
		q.Set(k, v)
	}
	endpoint.RawQuery = q.Encode()
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
	var v GetImagesResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	for k, i := range v.Images {
		v.Images[k].CreatedAt = toJst(i.CreatedAt)
		v.Images[k].UpdatedAt = toJst(i.UpdatedAt)
	}
	return &v, err
}

func (api *V3) GetUsedImageCapacity() {

}

func (api *V3) GetImageCapacity() {

}

func (api *V3) UpdateImageCapacity() {

}

func (api *V3) DeleteImage() {

}

func (api *V3) GetImage() {

}
