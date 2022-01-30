package amp360

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestGetListMock_withLogging(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"success":true,"message":"Successfully found the client's templates.","data":{"count":1,"rows":[{"id":1,"name":"APITEST","createdAt":"2021-11-18T06:17:45.000Z","updatedAt":"2021-11-18T06:17:45.000Z","ClientId":"ce16c215-e5a2-4ce6-9429-3bea82624a87","parentId":null,"Client":{"id":"ce16c215-e5a2-4ce6-9429-3bea82624a87","name":"TEST"},"Applications":[{"name":"TEST","version":"02.03.029","state":"Production","id":"766d0d8f-a0fd-4fa6-97e3-e44028305ba3","createdAt":"2021-11-10T06:15:53.000Z","fileName":"test.apk"}],"parentInfo":null}]}}`)
	})

	tl := TemplateList{}
	err := c.TemplatesService.GetList(context.Background(), nil, &tl)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 1

	if tl.Count != want {
		t.Errorf("Templates count = %v, want %v", tl.Count, want)
	}

}
