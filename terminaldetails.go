package amp360

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type Details struct {
	TemplateDetails []struct {
		ID            int       `json:"id"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
		AppTemplateID int       `json:"AppTemplateId"`
		ApplicationID string    `json:"ApplicationId"`
		AppTemplate   struct {
			ID        int       `json:"id"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"AppTemplate"`
		Application struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Version  string `json:"version"`
			State    string `json:"state"`
			Type     string `json:"type"`
			FileName string `json:"fileName"`
		} `json:"Application"`
	} `json:"templateDetails"`
	Terminal struct {
		ID              int         `json:"id"`
		SerialNumber    string      `json:"serialNumber"`
		Status          string      `json:"status"`
		Name            string      `json:"name"`
		Imei            interface{} `json:"imei"`
		EthernetMAC     interface{} `json:"ethernetMAC"`
		WifiMAC         interface{} `json:"wifiMAC"`
		BluetoothMAC    interface{} `json:"bluetoothMAC"`
		CloudAuthCode   string      `json:"cloudAuthCode"`
		QueueFirmware   int         `json:"queueFirmware"`
		CreatedAt       time.Time   `json:"createdAt"`
		UpdatedAt       time.Time   `json:"updatedAt"`
		AppTemplateID   int         `json:"AppTemplateId"`
		ClientID        string      `json:"ClientId"`
		FirmwareID      string      `json:"FirmwareId"`
		TerminalModelID string      `json:"TerminalModelId"`
		Firmware        struct {
			ID        string    `json:"id"`
			Name      string    `json:"name"`
			Version   string    `json:"version"`
			IsLatest  int       `json:"isLatest"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"Firmware"`
		TerminalModel struct {
			ID         string    `json:"id"`
			Name       string    `json:"name"`
			HardwareID string    `json:"hardwareId"`
			CreatedAt  time.Time `json:"createdAt"`
		} `json:"TerminalModel"`
	} `json:"terminal"`
}

func (c *TerminalsService) GetDetails(ctx context.Context, opt interface{}, v interface{}) (err error) {
	path := "terminals/details"
	var url *url.URL
	if url, err = addOptions(path, opt); err != nil {
		return err
	}

	return c.client.processRequest(ctx, http.MethodGet, *url, nil, v)
}
