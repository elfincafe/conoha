package v3

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/elfincafe/annette"
)

func (api *V3) PublishTokenById(uri, userId, password, tenantId string) (*annette.Response, error) {
	body := `
		{
			"auth": {
				"identity": {
					"methods": [
						"password"
					],
					"password": {
						"user": {
						"id": "###USER_ID###",
						"password": "###PASSWORD###"
						}
					}
				},
				"scope": {
					"project": {
						"id": "###TENANT_ID###"
					}
				}
			}
		}
  	`
	body = strings.ReplaceAll(body, "###USER_ID###", userId)
	body = strings.ReplaceAll(body, "###PASSWORD###", password)
	body = strings.ReplaceAll(body, "###TENANT_ID###", tenantId)
	return api.publishToken(uri, body)
}

func (api *V3) PublishTokenByName(uri, userName, password, tenantName string) (*annette.Response, error) {
	body := `
		{
			"auth": {
				"identity": {
					"methods": [
						"password"
					],
					"password": {
						"user": {
						"name": "###USER_NAME###",
						"password": "###PASSWORD###"
						}
					}
				},
				"scope": {
					"project": {
						"name": "###TENANT_NAME"
					}
				}
			}
		}
	`
	body = strings.ReplaceAll(body, "###USER_NAME###", userName)
	body = strings.ReplaceAll(body, "###PASSWORD###", password)
	body = strings.ReplaceAll(body, "###TENANT_NAME###", tenantName)
	return api.publishToken(uri, body)
}

func (api *V3) publishToken(uri, body string) (*annette.Response, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	client := annette.New(u)
	res, err := client.Post(body)
	if err != nil {
		return nil, err
	}
	api.Token = res.GetHeader("x-subject-token")
	// reading response body
	var jVal any
	err = json.Unmarshal(res.Binary(), &jVal)
	if err != nil {
		return nil, err
	}
	for k1, v1 := range jVal.(map[string]any) {
		if k1 != "token" {
			continue
		}
		for k2, v2 := range v1.(map[string]any) {
			if k2 == "issued_at" {
				api.IssuedAt, err = time.Parse(time.RFC3339Nano, v2.(string))
				continue
			} else if k2 == "expires_at" {
				api.ExpiredAt, err = time.Parse(time.RFC3339Nano, v2.(string))
				continue
			} else if k2 == "user" {
				for k3, v3 := range v2.(map[string]any) {
					if k3 == "id" {
						api.UserId = v3.(string)
					} else if k3 == "name" {
						api.UserName = v3.(string)
					}
				}
				continue
			} else if k2 == "project" {
				for k3, v3 := range v2.(map[string]any) {
					if k3 == "id" {
						api.TenantId = v3.(string)
					} else if k3 == "name" {
						api.TenantName = v3.(string)
					}
				}
				continue
			} else if k2 == "catalog" {
				for _, v3 := range v2.([]map[string]any) {
					for k4, v4 := range v3 {
						typ := ""
						if k4 == "type" {
							typ = v4.(string)
							continue
						} else if k4 != "endpoints" {
							continue
						}
						for k5, v5 := range v4.(map[string]string) {
							if k5 != "url" {
								continue
							}
							u, _ := url.Parse(v5)
							switch typ {
							case "identity":
								api.Endpoints.Identity = u
							case "compute":
								api.Endpoints.Compute = u
							case "load-balancer":
								api.Endpoints.LoadBalancer = u
							case "object-store":
								api.Endpoints.ObjectStorage = u
							case "dns":
								api.Endpoints.Dns = u
							case "volumev3":
								api.Endpoints.Volume = u
							case "image":
								api.Endpoints.Image = u
							case "network":
								api.Endpoints.Network = u
							case "account":
								api.Endpoints.Account = u
							}
						}
					}
				}
			}
		}
	}

	return res, nil
}
