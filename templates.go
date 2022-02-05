package amp360

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TemplatesService service

type TemplateList struct {
	Count int        `json:"count"`
	Rows  []Template `json:"rows"`
}

type Template struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	ClientID     string         `json:"ClientId"`
	ParentID     interface{}    `json:"parentId"`
	Client       TemplateClient `json:"Client"`
	Applications []AppTemplate  `json:"Applications"`
	ParentInfo   interface{}    `json:"parentInfo"`
}

type TemplateClient struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OriginPath string `json:"originPath,omitempty"`
}

type AppTemplate struct {
	Name      string    `json:"name"`
	Version   string    `json:"version,omitempty"`
	State     string    `json:"state,omitempty"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FileName  string    `json:"fileName,omitempty"`
}

func (c *TemplatesService) GetList(ctx context.Context, opt interface{}, v interface{}) (err error) {
	path := "templates"
	var url *url.URL
	if url, err = addOptions(path, opt); err != nil {
		return err
	}

	req, err := c.client.newRequestCtx(ctx, http.MethodGet, *url, nil)
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
