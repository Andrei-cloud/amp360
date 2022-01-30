package amp360

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestGetTerminalsListMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"success":true,"message":"Successfully found the terminals.","data":{"count":1,"rows":[{"id":25,"serialNumber":"8000044499","status":"Pending download","name":"Test Terminal 9","imei":null,"ethernetMAC":null,"wifiMAC":null,"bluetoothMAC":null,"cloudAuthCode":null,"queueFirmware":false,"createdAt":"2021-12-27T05:01:56.000Z","updatedAt":"2021-12-27T05:01:56.000Z","AppTemplateId":814,"ClientId":"test_client","FirmwareId":"test_firmware","TerminalModelId":"test1","AppTemplate":{"id":814,"name":"APITEST","createdAt":"2021-11-18T06:17:45.000Z"},"Client":{"id":"test_client","name":"TEST","originPath":"test"}}]}}`)
	})

	tl := TerminalsList{}
	err := c.TerminalsService.GetTerminalsList(context.Background(), nil, &tl)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 1

	if tl.Count != want {
		t.Errorf("Terminals count = %v, want %v", tl.Count, want)
	}

}

func TestCreateTerminal_NoTemplateMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading request body: %v", err)
		}

		values := NewTerminal{}
		if err := json.Unmarshal(b, &values); err != nil {
			t.Errorf("invalid body cannot parse %v", err)
		}
		want := NewTerminal{
			ModelID:      "value1",
			SerialNumber: "value2",
			Name:         "value3",
		}

		if !reflect.DeepEqual(values, want) {
			t.Errorf("invalid body received parsed as %v+", values)
		}
		fmt.Fprint(w, `{"id":321,"AppTemplateId":123,"ClientId":"test","FirmwareId":"test1","TerminalModelId":"test3","serialNumber":"80000123456","name":"test4","status":"Pending download","createdAt":"2022-01-30T15:30:36.441Z","updatedAt":"2022-01-30T15:30:36.441Z"}`)
	})

	nt := NewTerminal{
		ModelID:      "value1",
		SerialNumber: "value2",
		Name:         "value3",
	}
	ct := CreatedTerminal{}
	err := c.TerminalsService.CreateTerminal(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}

}

func TestCreateTerminal_WithTemplateMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading request body: %v", err)
		}

		values := NewTerminal{}
		if err := json.Unmarshal(b, &values); err != nil {
			t.Errorf("invalid body cannot parse %v", err)
		}
		want := NewTerminal{
			ModelID:      "value1",
			SerialNumber: "value2",
			Name:         "value3",
			TemplateID:   "123",
		}

		if !reflect.DeepEqual(values, want) {
			t.Errorf("invalid body received parsed as %v+", values)
		}
		fmt.Fprint(w, `{"id":321,"AppTemplateId":123,"ClientId":"test","FirmwareId":"test1","TerminalModelId":"test3","serialNumber":"80000123456","name":"test4","status":"Pending download","createdAt":"2022-01-30T15:30:36.441Z","updatedAt":"2022-01-30T15:30:36.441Z"}`)
	})

	nt := NewTerminal{
		ModelID:      "value1",
		SerialNumber: "value2",
		Name:         "value3",
		TemplateID:   "123",
	}
	ct := CreatedTerminal{}
	err := c.TerminalsService.CreateTerminal(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}

}

func TestCreateTerminal_WithClientMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading request body: %v", err)
		}

		values := NewTerminal{}
		if err := json.Unmarshal(b, &values); err != nil {
			t.Errorf("invalid body cannot parse %v", err)
		}
		want := NewTerminal{
			ModelID:      "value1",
			SerialNumber: "value2",
			Name:         "value3",
			ClientID:     "test",
		}

		if !reflect.DeepEqual(values, want) {
			t.Errorf("invalid body received parsed as %v+", values)
		}
		fmt.Fprint(w, `{"id":321,"AppTemplateId":123,"ClientId":"test","FirmwareId":"test1","TerminalModelId":"test3","serialNumber":"80000123456","name":"test4","status":"Pending download","createdAt":"2022-01-30T15:30:36.441Z","updatedAt":"2022-01-30T15:30:36.441Z"}`)
	})

	nt := NewTerminal{
		ModelID:      "value1",
		SerialNumber: "value2",
		Name:         "value3",
		ClientID:     "test",
	}
	ct := CreatedTerminal{}
	err := c.TerminalsService.CreateTerminal(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}
}

func TestCreateTerminal_WithParamsMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Content-Type", "application/json")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading request body: %v", err)
		}

		values := NewTerminal{}
		if err := json.Unmarshal(b, &values); err != nil {
			t.Errorf("invalid body cannot parse %v", err)
		}

		got := values.Parameters["key2"].(bool)
		if !got {
			t.Errorf("invalid body received parsed as %+v", got)
		}
		fmt.Fprint(w, `{"id":321,"AppTemplateId":123,"ClientId":"test","FirmwareId":"test1","TerminalModelId":"test3","serialNumber":"80000123456","name":"test4","status":"Pending download","createdAt":"2022-01-30T15:30:36.441Z","updatedAt":"2022-01-30T15:30:36.441Z"}`)
	})

	params := map[string]interface{}{
		"key1": "value1",
		"key2": true,
		"key3": "value3",
		"key4": 4,
	}
	nt := NewTerminal{
		ModelID:      "value1",
		SerialNumber: "value2",
		Name:         "value3",
		TemplateID:   "123",
		Parameters:   params,
	}
	ct := CreatedTerminal{}
	err := c.TerminalsService.CreateTerminal(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}
}

func BenchmarkGetTerminalsListMock(b *testing.B) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"message":"Successfully found the terminals.","data":{"count":1,"rows":[{"id":25,"serialNumber":"8000044499","status":"Pending download","name":"Test Terminal 9","imei":null,"ethernetMAC":null,"wifiMAC":null,"bluetoothMAC":null,"cloudAuthCode":null,"queueFirmware":false,"createdAt":"2021-12-27T05:01:56.000Z","updatedAt":"2021-12-27T05:01:56.000Z","AppTemplateId":814,"ClientId":"test_client","FirmwareId":"test_firmware","TerminalModelId":"test1","AppTemplate":{"id":814,"name":"APITEST","createdAt":"2021-11-18T06:17:45.000Z"},"Client":{"id":"test_client","name":"TEST","originPath":"test"}}]}}`)
	})

	tl := TerminalsList{}
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		err := c.TerminalsService.GetTerminalsList(context.Background(), nil, &tl)
		if err != nil {
			b.Errorf("Error occured = %v", err)
		}
	}
}
