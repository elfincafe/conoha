package conoha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elfincafe/annette"
	"github.com/google/uuid"
)

type (
	Domain struct {
		Uuid      uuid.UUID `json:"uuid"`
		Name      string    `json:"name"`
		ProjectId string    `json:"project_id"`
		Serial    int       `json:"serial"`
		Ttl       int       `json:"ttl"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	Record struct {
		Uuid       uuid.UUID `json:"uuid"`
		DomainUuid uuid.UUID `json:"domain_uuid"`
		Name       string    `json:"name"`
		Type       string    `json:"type"`
		Data       string    `json:"data"`
		Priority   int       `json:"priority"`
		Weight     int       `json:"weight"`
		Port       int       `json:"port"`
		Ttl        int       `json:"ttl"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
	recordRequest struct {
		Name     string `json:"name,omitempty"`
		Type     string `json:"type,omitempty"`
		Data     string `json:"data,omitempty"`
		Priority string `json:"priority,omitempty"`
		Weight   string `json:"weight,omitempty"`
		Port     string `json:"port,omitempty"`
	}
	GetDomainsResponse struct {
		Domains    []Domain `json:"domains"`
		TotalCount int      `json:"total_count"`
	}
	UpdateDomainResponse Domain
	CreateDomainResponse Domain
	GetDomainResponse    Domain
	GetRecordsResponse   struct {
		Records    []Record `json:"records"`
		TotalCount int      `json:"total_count"`
	}
	CreateRecordResponse Record
	UpdateRecordResponse Record
	GetRecordResponse    Record
)

func (api *V3) GetDomains(limit, offset int, sort, key string) (*GetDomainsResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = "/v1/domains"
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	sort = strings.ToLower(sort)
	if sort != "desc" {
		sort = "asc"
	}
	key = strings.ToLower(key)
	switch key {
	case "uuid", "name", "project_id", "serial", "email", "created_at", "updated_at":
		// Do Nothing
	default:
		key = "created_at"
	}
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
	var v GetDomainsResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	for k, d := range v.Domains {
		v.Domains[k].CreatedAt = toJst(d.CreatedAt)
		v.Domains[k].UpdatedAt = toJst(d.UpdatedAt)
	}
	return &v, nil
}

func (api *V3) DeleteDomain(domainId uuid.UUID) error {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf("/v1/domains/%s", domainId)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Delete()
	if err != nil {
		return err
	}
	if !res.IsStatus204() {
		return toError(res.Binary())
	}
	return nil
}

func (api *V3) UpdateDomain(domainId uuid.UUID, email string, ttl int) (*UpdateDomainResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf("/v1/domains/%s", domainId)
	body := fmt.Sprintf(`{
		"ttl": %d,
		"email": "%s"
	}`, ttl, email)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Put(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v UpdateDomainResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
}

func (api *V3) CreateDomain(domain, email string, ttl int) (*CreateDomainResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = "/v1/domains"
	domain = strings.Trim(domain, "\r\n\t\v .") + "."
	body := fmt.Sprintf(`{
		"name":"%s",
		"ttl": %d,
		"email": "%s"
	}`, domain, ttl, email)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v CreateDomainResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
}

func (api *V3) GetDomain(domainId uuid.UUID) (*GetDomainResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf(`/v1/domains/%s`, domainId)
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
	var v GetDomainResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
}

func (api *V3) GetRecords(domainId uuid.UUID, limit, offset int, sort, key string) (*GetRecordsResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf(`/v1/domains/%s/records`, domainId)
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	sort = strings.ToLower(sort)
	if sort != "desc" {
		sort = "asc"
	}
	key = strings.ToLower(key)
	switch key {
	case "uuid", "name", "project_id", "serial", "email", "created_at", "updated_at":
		// Do Nothing
	default:
		key = "created_at"
	}
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
	var v GetRecordsResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	for k, r := range v.Records {
		v.Records[k].CreatedAt = toJst(r.CreatedAt)
		v.Records[k].UpdatedAt = toJst(r.UpdatedAt)
	}
	return &v, nil
}

func (api *V3) CreateRecord(domainId uuid.UUID, name, recType, data, priority, weight, port string) (*CreateRecordResponse, error) {
	req := recordRequest{}
	req.Name = strings.Trim(name, "\r\n\t\v .") + "."
	req.Data = data
	recType = strings.ToUpper(recType)
	switch recType {
	case "A", "AAAA", "CNAME", "NS", "TXT":
		req.Type = recType
	case "MX":
		req.Type = recType
		req.Priority = priority
	case "SRV":
		req.Type = recType
		req.Priority = priority
		req.Weight = weight
		req.Port = port
	default:
		req.Type = recType
	}
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf(`/v1/domains/%s/records`, domainId)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Post(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v CreateRecordResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
}

func (api *V3) DeleteRecord(domainId, recordId uuid.UUID) error {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf("/v1/domains/%s/records/%s", domainId, recordId)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)
	res, err := client.Delete()
	if err != nil {
		return err
	}
	if !res.IsStatus204() {
		return toError(res.Binary())
	}
	return nil
}

func (api *V3) UpdateRecord(domainId, recordId uuid.UUID, name, recType, data, priority, weight, port string) (*UpdateRecordResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf("/v1/domains/%s/records/%s", domainId, recordId)
	client := annette.New(endpoint)
	client.Header.Set("Accept", "application/json")
	client.Header.Set("Content-Type", "application/json")
	client.Header.Set("X-Auth-Token", api.Token)

	// request data
	req := recordRequest{}
	name = strings.Trim(name, "\r\n\t\v .")
	if name != "" {
		name += "."
	}
	req.Name = name
	req.Data = data
	recType = strings.ToUpper(recType)
	switch recType {
	case "A", "AAAA", "CNAME", "NS", "TXT":
		req.Type = recType
	case "MX":
		req.Type = recType
		req.Priority = priority
	case "SRV":
		req.Type = recType
		req.Priority = priority
		req.Weight = weight
		req.Port = port
	default:
		req.Type = recType
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := client.Put(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if !res.IsStatus200() {
		return nil, toError(res.Binary())
	}
	var v UpdateRecordResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
}

func (api *V3) GetRecord(domainId, recordId uuid.UUID) (*GetRecordResponse, error) {
	endpoint := api.Endpoints.Dns
	endpoint.Path = fmt.Sprintf(`/v1/domains/%s/records/%s`, domainId, recordId)
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
	var v GetRecordResponse
	err = json.Unmarshal(res.Binary(), &v)
	if err != nil {
		return nil, err
	}
	v.CreatedAt = toJst(v.CreatedAt)
	v.UpdatedAt = toJst(v.UpdatedAt)
	return &v, nil
}
