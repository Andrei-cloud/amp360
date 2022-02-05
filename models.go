package amp360

import (
	"context"
	"net/http"
	"net/url"
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

func (c *ModelsService) GetList(ctx context.Context, v interface{}) (err error) {
	path := "models"
	url := url.URL{Path: path}
	return c.client.processRequest(ctx, http.MethodGet, url, nil, v)
}
