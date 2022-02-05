package amp360

import (
	"context"
	"net/http"
	"net/url"
)

type CompaniesService service

type CompaniesList struct {
	Count int       `json:"count"`
	Rows  []Company `json:"rows"`
}

type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CompaniesOpt struct {
	Size int `url:"size"`
	Page int `url:"page"`
}

func (c *CompaniesService) GetList(ctx context.Context, opt interface{}, v interface{}) (err error) {
	path := "client/children"
	var url *url.URL
	if url, err = addOptions(path, opt); err != nil {
		return err
	}

	return c.client.processRequest(ctx, http.MethodGet, *url, nil, v)
}
