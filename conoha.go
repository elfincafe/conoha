package v3

import (
	"net/url"
	"time"
)

type (
	Conoha interface {
		Version() string
	}
	V3 struct {
		UserId     string    `json:"user_id"`
		UserName   string    `json:"user_name"`
		TenantId   string    `json:"tenant_id"`
		TenantName string    `json:"tenant_name"`
		Token      string    `json:"token"`
		IssuedAt   time.Time `json:"issued_at"`
		ExpiredAt  time.Time `json:"expires_at"`
		Endpoints  Endpoint  `json:"endpoints"`
	}
	V2 struct {
		UserId     string    `json:"user_id"`
		UserName   string    `json:"user_name"`
		TenantId   string    `json:"tenant_id"`
		TenantName string    `json:"tenant_name"`
		Token      string    `json:"token"`
		IssuedAt   time.Time `json:"issued_at"`
		ExpiredAt  time.Time `json:"expires_at"`
		Endpoints  Endpoint  `json:"endpoints"`
	}
	Endpoint struct {
		Identity      *url.URL `json:"identity,omitempty"`
		Compute       *url.URL `json:"compute,omitempty"`
		LoadBalancer  *url.URL `json:"load_balancer,omitempty"`
		ObjectStorage *url.URL `json:"object_storage,omitempty"`
		Dns           *url.URL `json:"dns,omitempty"`
		Volume        *url.URL `json:"volume,omitempty"`
		Image         *url.URL `json:"image,omitempty"`
		Network       *url.URL `json:"network,omitempty"`
		Account       *url.URL `json:"account,omitempty"`
		S3            *url.URL `json:"s3,omitempty"`
		Database      *url.URL `json:"url,omitempty"`
	}
)

func New(ver string) Conoha {
	if ver == "v3" {
		return &V3{
			IssuedAt:  time.Date(1970, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
			ExpiredAt: time.Date(1970, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
			Endpoints: Endpoint{},
		}
	} else if ver == "v2" {
		return &V2{
			IssuedAt:  time.Date(1970, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
			ExpiredAt: time.Date(1970, 1, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*60*60)),
			Endpoints: Endpoint{},
		}
	} else {
		return nil
	}
}

func (api *V3) Version() string {
	return "v3"
}

func (api *V2) Version() string {
	return "v2"
}
