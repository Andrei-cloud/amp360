package amp360

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestModelsGetListMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/models", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"success":true,"message":"Successfully found available terminal models.","data":{"count":3,"rows":[{"name":"TEST1","id":"test1","hardwareId":"CD","maintenanceInterval":180,"jointName":"TEST1-CD","createdAt":"2021-10-30T00:55:39.000Z"},{"name":"TEST2","id":"test2","hardwareId":"2AA","maintenanceInterval":180,"jointName":"TEST2-2AA","createdAt":"2021-10-30T00:55:39.000Z"},{"name":"TEST3","id":"test3","hardwareId":"2AA","maintenanceInterval":180,"jointName":"TEST3-2AA","createdAt":"2021-10-30T00:55:39.000Z"}]}}`)
	})

	ml := ModelsList{}
	err := c.ModelsService.GetList(context.Background(), &ml)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 3

	if ml.Count != want {
		t.Errorf("Models count = %v, want %v", ml.Count, want)
	}

}

func BenchmarkModelsGetListMock(b *testing.B) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/models", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"message":"Successfully found available terminal models.","data":{"count":3,"rows":[{"name":"TEST1","id":"test1","hardwareId":"CD","maintenanceInterval":180,"jointName":"TEST1-CD","createdAt":"2021-10-30T00:55:39.000Z"},{"name":"TEST2","id":"test2","hardwareId":"2AA","maintenanceInterval":180,"jointName":"TEST2-2AA","createdAt":"2021-10-30T00:55:39.000Z"},{"name":"TEST3","id":"test3","hardwareId":"2AA","maintenanceInterval":180,"jointName":"TEST3-2AA","createdAt":"2021-10-30T00:55:39.000Z"}]}}`)
	})

	ml := ModelsList{}
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		err := c.ModelsService.GetList(context.Background(), &ml)
		if err != nil {
			b.Errorf("Error occured = %v", err)
		}
	}
}
