package amp360

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *TerminalsService) GetParams(ctx context.Context, templateID string, opt interface{}, v interface{}) (err error) {
	if templateID == "" {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/params/%s", templateID)
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
