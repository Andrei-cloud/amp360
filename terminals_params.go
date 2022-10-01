package amp360

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Parameter struct {
	ID                int         `json:"id"`
	Type              string      `json:"type"`
	Tag               string      `json:"tag"`
	Name              string      `json:"name"`
	Hint              string      `json:"hint"`
	Validator         string      `json:"validator"`
	Value             string      `json:"value"`
	DefaultValue      string      `json:"defaultValue"`
	VisibleOnTemplate int         `json:"visibleOnTemplate"`
	VisibleOnTerminal int         `json:"visibleOnTerminal"`
	FilePath          interface{} `json:"filePath"`
	ApplicationID     string      `json:"ApplicationId"`
	ParamCategoryID   string      `json:"ParamCategoryId"`
	CategoryName      string      `json:"categoryName"`
}
type TerminalParams struct {
	Categories []Categories `json:"categories"`
	Count      int          `json:"count"`
	Rows       []Parameter  `json:"rows"`
}

func (c *TerminalsService) GetParams(ctx context.Context, id int, opt interface{}, v interface{}) (err error) {
	if id == 0 {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/params/%d", id)

	var url *url.URL
	if url, err = addOptions(path, opt); err != nil {
		return err
	}

	return c.client.processRequest(ctx, http.MethodGet, *url, nil, v)
}

func (c *TerminalsService) UpdateParams(ctx context.Context, id int, params map[string]string, paramfiles map[string]string, u, f interface{}) (err error) {
	if id == 0 {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/params/bulk/%d", id)
	url := url.URL{Path: path}
	return c.client.processBulkRequest(ctx, http.MethodPost, url, params, paramfiles, u, f)
}
