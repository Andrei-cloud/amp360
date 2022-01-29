package amp360

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TerminalsService service

type Terminal struct {
	ID              int         `json:"id"`           // used in create response
	SerialNumber    string      `json:"serialNumber"` // used in create response
	Status          string      `json:"status"`       // used in create response
	Name            string      `json:"name"`         // used in create response
	Imei            string      `json:"imei,omitempty"`
	EthernetMAC     string      `json:"ethernetMAC"`
	WifiMAC         string      `json:"wifiMAC,omitempty"`
	BluetoothMAC    string      `json:"bluetoothMAC,omitempty"`
	CloudAuthCode   interface{} `json:"cloudAuthCode,omitempty"`
	QueueFirmware   bool        `json:"queueFirmware"`   // used in create response
	CreatedAt       time.Time   `json:"createdAt"`       // used in create response
	UpdatedAt       time.Time   `json:"updatedAt"`       // used in create response
	AppTemplateID   int         `json:"AppTemplateId"`   // used in create response
	ClientID        string      `json:"ClientId"`        // used in create response
	FirmwareID      string      `json:"FirmwareId"`      // used in create response
	TerminalModelID string      `json:"TerminalModelId"` // used in create response
	AppTemplate     struct {
		Name      string    `json:"name"`
		ID        int       `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"AppTemplate,omitempty"`
	Client Client `json:"Client,omitempty"`
}

type TerminalsList struct {
	Count int        `json:"count"`
	Rows  []Terminal `json:"rows"`
}

type TerminalsOpt struct {
	Size int `url:"size"`
	Page int `url:"page"`
}

func (c *TerminalsService) GetTerminalsList(ctx context.Context, opt interface{}, v interface{}) (err error) {
	path := "terminals"
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
