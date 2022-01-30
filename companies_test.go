package amp360

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestCompaniesGetListMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/client/children", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"success":true,"message":"Successfully fetched sub-clients.","data":{"count":2,"rows":[{"id":"test1","name":"TEST 1","type":"MERCHANT"},{"id":"test2","name":"TEST2","type":"MERCHANT"}]}}`)
	})

	cl := CompaniesList{}
	err := c.CompaniesService.GetList(context.Background(), nil, &cl)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 2

	if cl.Count != want {
		t.Errorf("Companies count = %v, want %v", cl.Count, want)
	}

}

func BenchmarkCompaniesGetListMock(b *testing.B) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/client/children", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"message":"Successfully fetched sub-clients.","data":{"count":2,"rows":[{"id":"test1","name":"TEST 1","type":"MERCHANT"},{"id":"test2","name":"TEST2","type":"MERCHANT"}]}}`)
	})

	cl := CompaniesList{}
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		err := c.CompaniesService.GetList(context.Background(), nil, &cl)
		if err != nil {
			b.Errorf("Error occured = %v", err)
		}
	}
}
