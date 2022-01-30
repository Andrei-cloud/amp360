package amp360

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type TemplateParams struct {
	Categories []Categories `json:"categories"`
	Count      int          `json:"count"`
	Rows       []Param      `json:"rows"`
}
type Categories struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Param struct {
	ID                 int       `json:"id"`
	Type               string    `json:"type"`
	Tag                string    `json:"tag"`
	Name               string    `json:"name"`
	Hint               string    `json:"hint"`
	Validator          string    `json:"validator"`
	VisibleOnTemplate  int       `json:"visibleOnTemplate"`
	VisibleOnTerminal  int       `json:"visibleOnTerminal"`
	EditableOnTerminal int       `json:"editableOnTerminal"`
	FilePath           string    `json:"filePath"`
	ApplicationID      string    `json:"ApplicationId"`
	ParamCategoryID    string    `json:"ParamCategoryId"`
	CategoryName       string    `json:"categoryName"`
	Value              string    `json:"value"`
	DefaultValue       string    `json:"defaultValue"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type ParamsOpt struct {
	CategoryId string `url:"categoryId"`
}

func (c *TemplatesService) GetParams(ctx context.Context, templateID string, opt interface{}, v interface{}) (err error) {
	if templateID == "" {
		return errors.New("required templateID is missing")
	}
	path := fmt.Sprintf("templates/params/%s", templateID)
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
	if res.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

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
