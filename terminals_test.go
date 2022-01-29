package amp360

import (
	"context"
	"fmt"
	"net/http"
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
