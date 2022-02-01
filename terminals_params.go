package amp360

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func (c *TerminalsService) UpdateParams(ctx context.Context, templateID string, params map[string]string, paramfiles map[string]string, u, f interface{}) (err error) {
	if templateID == "" {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/params/%s", templateID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for p, filePath := range paramfiles {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		part, err := writer.CreateFormFile(p, filepath.Base(file.Name()))
		if err != nil {
			return err
		}
		io.Copy(part, file)
		if err != nil {
			return err
		}
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	writer.Close()
	req, err := c.client.NewMultiPartRequestCtx(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	resp := BulkResponse{
		Failed:  f,
		Updated: u,
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
