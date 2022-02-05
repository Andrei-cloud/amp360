package amp360

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

var (
	terminalsRe = regexp.MustCompile(`^\/terminals\/(\d+)`)
)

func TestDeleteMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !terminalsRe.MatchString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			testMethod(t, r, http.MethodDelete)
			fmt.Fprint(w, `{"success":true,"message":"Successfully deleted the terminal.","data":{}}`)
		}
	})

	wantErr := errors.New("required terminalID is missing")
	err := c.TerminalsService.Delete(context.Background(), 0)
	if err == nil {
		t.Errorf("Error is nil, want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("Error got %v, want %v", err, wantErr)
	}

	err = c.TerminalsService.Delete(context.Background(), 321)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}
}

func TestUpdateMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !terminalsRe.MatchString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			testMethod(t, r, http.MethodPut)
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
				Name:     "New name",
				ClientID: "new client",
			}

			if !reflect.DeepEqual(values, want) {
				t.Errorf("invalid body received parsed as %v+", values)
			}
			fmt.Fprint(w, `{"success":true,"message":"Successfully updated the terminal.","data":{}}`)
		}
	})

	ut := NewTerminal{
		Name:     "New name",
		ClientID: "new client",
	}

	wantErr := errors.New("required terminalID is missing")
	err := c.TerminalsService.Update(context.Background(), 0, &ut)
	if err == nil {
		t.Errorf("Error is nil, want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("Error got %v, want %v", err, wantErr)
	}

	err = c.TerminalsService.Update(context.Background(), 321, &ut)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}
}

func TestGetListMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"success":true,"message":"Successfully found the terminals.","data":{"count":1,"rows":[{"id":25,"serialNumber":"8000044499","status":"Pending download","name":"Test Terminal 9","imei":null,"ethernetMAC":null,"wifiMAC":null,"bluetoothMAC":null,"cloudAuthCode":null,"queueFirmware":false,"createdAt":"2021-12-27T05:01:56.000Z","updatedAt":"2021-12-27T05:01:56.000Z","AppTemplateId":814,"ClientId":"test_client","FirmwareId":"test_firmware","TerminalModelId":"test1","AppTemplate":{"id":814,"name":"APITEST","createdAt":"2021-11-18T06:17:45.000Z"},"Client":{"id":"test_client","name":"TEST","originPath":"test"}}]}}`)
	})

	tl := TerminalsList{}
	err := c.TerminalsService.GetList(context.Background(), nil, &tl)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 1

	if tl.Count != want {
		t.Errorf("Terminals count = %v, want %v", tl.Count, want)
	}

}

func TestGetDetailsMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/terminals/details", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"success":true,"message":"Successfully found the terminal details.","data":{"templateDetails":[{"id":393,"createdAt":"2021-10-30T04:29:58.000Z","updatedAt":"2021-10-30T04:29:58.000Z","AppTemplateId":234,"ApplicationId":"6fa5752f","AppTemplate":{"id":234,"name":"TEST-Template-22","createdAt":"2021-10-30T04:29:58.000Z"},"Application":{"id":"6fa5752f","name":"TEST","version":"v02.02.022","state":"Testing","type":"Payment","fileName":"AMP POS v02.02.022_TESTta_01.00.001.apk"}}],"terminal":{"id":151,"serialNumber":"8000000789","status":"Migrated","name":"98765432","imei":null,"ethernetMAC":null,"wifiMAC":null,"bluetoothMAC":null,"cloudAuthCode":"","queueFirmware":0,"createdAt":"2021-10-30T04:29:58.000Z","updatedAt":"2021-10-30T04:29:58.000Z","AppTemplateId":234,"ClientId":"07d018a5","FirmwareId":"2a0f80ef","TerminalModelId":"10af3b83","Firmware":{"id":"2a0f80ef","name":"AMP8000-2AA","version":"03.02.39","isLatest":1,"createdAt":"2021-10-30T00:55:39.000Z"},"TerminalModel":{"id":"10af3b83","name":"AMP8000","hardwareId":"2AA","createdAt":"2021-10-30T00:55:39.000Z"}}}}`)
	})

	term := Details{}
	opt := &TerminalsOpt{
		ID:           0,
		SerialNumber: "80002456",
		TID:          "",
		MID:          "",
	}
	err := c.TerminalsService.GetDetails(context.Background(), opt, &term)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 151

	if term.Terminal.ID != want {
		t.Errorf("Terminals ID = %v, want %v", term.Terminal.ID, want)
	}

}
func TestCreate_NoTemplateMock(t *testing.T) {
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
	err := c.TerminalsService.Create(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}

}

func TestCreate_WithTemplateMock(t *testing.T) {
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
	err := c.TerminalsService.Create(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}

}

func TestCreate_WithClientMock(t *testing.T) {
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
	err := c.TerminalsService.Create(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}
}

func TestCreate_WithParamsMock(t *testing.T) {
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
	err := c.TerminalsService.Create(context.Background(), &nt, &ct)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := "Pending download"

	if ct.Status != want {
		t.Errorf("Terminal status = %v, want %v", ct.Status, want)
	}
}

func BenchmarkGetListMock(b *testing.B) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/terminals", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"message":"Successfully found the terminals.","data":{"count":1,"rows":[{"id":25,"serialNumber":"8000044499","status":"Pending download","name":"Test Terminal 9","imei":null,"ethernetMAC":null,"wifiMAC":null,"bluetoothMAC":null,"cloudAuthCode":null,"queueFirmware":false,"createdAt":"2021-12-27T05:01:56.000Z","updatedAt":"2021-12-27T05:01:56.000Z","AppTemplateId":814,"ClientId":"test_client","FirmwareId":"test_firmware","TerminalModelId":"test1","AppTemplate":{"id":814,"name":"APITEST","createdAt":"2021-11-18T06:17:45.000Z"},"Client":{"id":"test_client","name":"TEST","originPath":"test"}}]}}`)
	})

	tl := TerminalsList{}
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		err := c.TerminalsService.GetList(context.Background(), nil, &tl)
		if err != nil {
			b.Errorf("Error occured = %v", err)
		}
	}
}
