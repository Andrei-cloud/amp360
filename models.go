package amp360

import (
	"context"
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

func (c *ModelsService) GetList(ctx context.Context, v interface{}) (err error) {
	path := "models"

	return c.client.processRequest(ctx, http.MethodGet, path, nil, v)
}
