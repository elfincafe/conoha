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
	image struct {
		HwRescueBus            string    `json:"hw_rescue_bus,omitempty"`
		HwRescueDevice         string    `json:"hw_rescue_device,omitempty"`
		HwVifMultiqueueEnabled bool      `json:"hw_vif_multiqueue_enabled,omitempty"`
		HwQemuGuestAgent       bool      `json:"hw_qemu_guest_agent,omitempty"`
		HwVideoModel           string    `json:"hw_video_model,omitempty"`
		Architecture           string    `json:"architecture,omitempty"`
		Bootable               bool      `json:"bootable"`
		Name                   string    `json:"name"`
		DiskFormat             string    `json:"disk_format"`
		ContainerFormat        string    `json:"container_format"`
		Visibility             string    `json:"visibility"`
		Size                   int       `json:"size"`
		VirtualSize            int       `json:"virtual_size"`
		Status                 string    `json:"status"`
		Checksum               int       `json:"checksum"`
		Protected              bool      `json:"protected"`
		MinRam                 int       `json:"min_ram"`
		MinDisk                int       `json:"min_disk"`
		Owner                  string    `json:"owner"`
		OsType                 string    `json:"os_type,omitempty"`
		OsHidden               bool      `json:"os_hidden"`
		OsHashAlgo             string    `json:"os_hash_algo"`
		OsHashValue            string    `json:"os_hash_value"`
		Id                     uuid.UUID `json:"id"`
		CreatedAt              time.Time `json:"created_at"`
		UpdatedAt              time.Time `json:"updated_at"`
		Tags                   []string  `json:"tags"`
		Self                   string    `json:"self"`
		File                   string    `json:"file"`
		Schema                 string    `json:"schema"`
	}
	CreateIsoImageResponse image
	GetImagesResponse      struct {
		Images []image `json:"images"`
		Schema string  `json:"schema"`
		First  string  `json:"first"`
	}
	GetUsedImageCapacityResponse struct {
		Images []struct {
			Size int
		} ``
	}
	GetImageCapacityResponse struct {
		Quota []struct {
			ImageSize string `json:"image_size"`
		} `json:"quota"`
	}
	UpdateImageCapacityResponse GetImageCapacityResponse
	GetImageResponse            image
)

func (api *V3) UploadIsoImage(imageId uuid.UUID, path string) error {
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
		return toError(res.Binary())
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
		return nil, toError(res.Binary())
	}
	var v CreateIsoImageResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
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
		return nil, toError(res.Binary())
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
	return &v, nil
}

func (api *V3) GetUsedImageCapacity() (*GetUsedImageCapacityResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = "/v2/images/total"
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Get()
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v GetUsedImageCapacityResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (api *V3) GetImageCapacity() (*GetImageCapacityResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = "/v2/quota"
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Get()
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v GetImageCapacityResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (api *V3) UpdateImageCapacity(imageSize string) (*UpdateImageCapacityResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = "/v2/quota"
	body := fmt.Sprintf(`{
		"quota": {"image_size": "%s"}
	}`, imageSize)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Put(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v UpdateImageCapacityResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (api *V3) DeleteImage(imageId uuid.UUID) error {
	endpoint := api.Endpoints.Image
	endpoint.Path = fmt.Sprintf("/v2/images/%s", imageId)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Delete()
	if err != nil {
		return err
	}
	if !res.IsStatus200() {
		return toError(res.Binary())
	}
	return nil
}

func (api *V3) GetImage(imageId uuid.UUID) (*GetImageResponse, error) {
	endpoint := api.Endpoints.Image
	endpoint.Path = fmt.Sprintf("/v2/images/%s", imageId)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Get()
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v GetImageResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
