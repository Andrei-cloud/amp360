package amp360

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func (c *TerminalsService) GetParams(ctx context.Context, templateID string, opt interface{}, v interface{}) (err error) {
	if templateID == "" {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/params/%s", templateID)
	var url *url.URL
	if url, err = addOptions(path, opt); err != nil {
		return err
	}

	return c.client.processRequest(ctx, http.MethodGet, *url, nil, v)
}

func (c *TerminalsService) UpdateParams(ctx context.Context, templateID string, params map[string]string, paramfiles map[string]string, u, f interface{}) (err error) {
	if templateID == "" {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/params/%s", templateID)
	url := url.URL{Path: path}
	return c.client.processBulkRequest(ctx, http.MethodPost, url, params, paramfiles, u, f)
}
