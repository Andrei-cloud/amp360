package amp360

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ModelsService service

type ModelsList struct {
	Count int             `json:"count"`
	Rows  []TerminalModel `json:"rows"`
}

type TerminalModel struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	HardwareID          string    `json:"hardwareId"`
	JointName           string    `json:"jointName"`
	MaintenanceInterval int       `json:"maintenanceInterval"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

func (c *ModelsService) GetModelsList(ctx context.Context, v interface{}) (err error) {
	path := "models"

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
