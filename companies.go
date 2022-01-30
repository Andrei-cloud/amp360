package amp360

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	if path, err = addOptions(path, opt); err != nil {
		return err
	}

	req, err := c.client.NewRequestCtx(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resp := Response{
		Data: v,
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if !resp.Success {
		err = fmt.Errorf("api err: %s", resp.Message)
	}

	return err
}
