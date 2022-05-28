package amp360

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	ID           int    `url:"id,omitempty"`
	SerialNumber string `url:"serialNumber,omitempty"`
	TID          string `url:"tid,omitempty"`
	MID          string `url:"mid,omitempty"`
	Size         int    `url:"size,omitempty"`
	Page         int    `url:"page,omitempty"`
}

type NewTerminal struct {
	ModelID      interface{}            `json:"modelId"`
	SerialNumber string                 `json:"serialNumber"`
	Name         string                 `json:"name"`
	ClientID     interface{}            `json:"clientId,omitempty"`
	TemplateID   string                 `json:"templateId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters"`
}

type CreatedTerminal struct {
	ID              int       `json:"id"`
	AppTemplateID   int       `json:"AppTemplateId"`
	ClientID        string    `json:"ClientId"`
	FirmwareID      string    `json:"FirmwareId"`
	TerminalModelID string    `json:"TerminalModelId"`
	SerialNumber    string    `json:"serialNumber"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (c *TerminalsService) GetList(ctx context.Context, opt interface{}, v interface{}) (err error) {
	path := "terminals"
	var url *url.URL
	if url, err = addOptions(path, opt); err != nil {
		return err
	}

	return c.client.processRequest(ctx, http.MethodGet, *url, nil, v)
}

func (c *TerminalsService) Create(ctx context.Context, data *NewTerminal, v interface{}) (err error) {
	path := "terminals"
	rel := &url.URL{Path: path}
	if data == nil {
		return errors.New("can't create terminals on nil data")
	}

	req, err := c.client.newRequestCtx(ctx, http.MethodPost, *rel, data)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK && res.StatusCode > 299 {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}

		resp := Response{}
		err = json.NewDecoder(res.Body).Decode(&resp)
		if err != nil {
			return err
		}

		if !resp.Success {
			err = fmt.Errorf("api err: %s", resp.Message)
		}

	} else {
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *TerminalsService) Update(ctx context.Context, id int, data *NewTerminal) (err error) {
	if id == 0 {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/%d", id)

	if data == nil {
		return errors.New("can't update terminal on nil data")
	}
	url := url.URL{Path: path}
	return c.client.processRequest(ctx, http.MethodPut, url, data, nil)
}

func (c *TerminalsService) Delete(ctx context.Context, id int) (err error) {
	if id == 0 {
		return errors.New("required terminalID is missing")
	}
	path := fmt.Sprintf("terminals/%d", id)
	url := url.URL{Path: path}
	return c.client.processRequest(ctx, http.MethodDelete, url, nil, nil)
}
