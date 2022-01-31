package amp360

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"
)

var (
	templateRe      = regexp.MustCompile(`^\/templates\/params\/(\d+)`)
	templateReQuery = regexp.MustCompile(`^\/templates\/params\/(\d+)\?categoryId=([a-z0-9\-]+)\&?`)
)

func TestTemplatesGetParamsMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !templateRe.MatchString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"success":true,"message":"Successfully found the template parameters.","data":{"categories":[{"name":"AMP Cloud","id":"6ba90b8c-0340-44fe-82c4-4bc32e51a316"},{"name":"Communications","id":"ca8958a4-8c21-4f82-af4a-58a66ac36e4a"}],"count":2,"rows":[{"id":950024,"type":"STRING","tag":"CLOUD.AUTHCODE","name":"CLOUD.AUTHCODE","hint":"","validator":"","value":"testtoken","defaultValue":"testtoken","visibleOnTemplate":1,"visibleOnTerminal":0,"filePath":"","ApplicationId":"763d0d8f-a0fd-4fa6-97e3-e44028305ba3","ParamCategoryId":"6ba90b8c-0340-44fe-82c4-4bc32e51a316","categoryName":"AMP Cloud","updatedAt":"2021-11-18T06:18:22.000Z","createdAt":"2021-11-18T06:18:22.000Z"},{"id":950047,"type":"STRING","tag":"COMMUNICATIONS.MEDIA.PRIMARY","name":"COMMUNICATIONS.MEDIA.PRIMARY","hint":"","validator":"","value":"CELLULAR","defaultValue":"CELLULAR","visibleOnTemplate":1,"visibleOnTerminal":0,"filePath":"","ApplicationId":"763d0d8f-a0fd-4fa6-97e3-e44028305ba3","ParamCategoryId":"ca8958a4-8c21-4f82-af4a-58a66ac36e4a","categoryName":"Communications","updatedAt":"2021-11-18T06:18:22.000Z","createdAt":"2021-11-18T06:18:22.000Z"}]}}`)
		}
	})

	tp := TemplateParams{}
	wantErr := errors.New("required templateID is missing")
	err := c.TemplatesService.GetParams(context.Background(), "", nil, &tp)
	if err == nil {
		t.Errorf("Error is nil, want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("Error got %v, want %v", err, wantErr)
	}

	want := 2
	err = c.TemplatesService.GetParams(context.Background(), "814", nil, &tp)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}
	if len(tp.Categories) != want {
		t.Errorf("Categorues count got %v, wnat %v", len(tp.Categories), want)
	}
	if tp.Count != want {
		t.Errorf("Parameters count got %v, wnat %v", len(tp.Categories), want)
	}
	if len(tp.Rows) != want {
		t.Errorf("Parameters actual count got %v, wnat %v", len(tp.Rows), want)
	}
}

func TestTemplatesGetParamsQueryMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !templateReQuery.MatchString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			url, _ := r.URL.Parse(r.URL.Path)
			testMethod(t, r, http.MethodGet)
			values := url.Query().Get("categoryId")
			want := "value1"
			if values != want {
				t.Errorf("invalid query received %v, want %v", values, want)
			}
			fmt.Fprint(w, `{"success":true,"message":"Successfully found the template parameters.","data":{"categories":[{"name":"AMP Cloud","id":"6ba90b8c-0340-44fe-82c4-4bc32e51a316"}],"count":1,"rows":[{"id":950024,"type":"STRING","tag":"CLOUD.AUTHCODE","name":"CLOUD.AUTHCODE","hint":"","validator":"","value":"testtoken","defaultValue":"testtoken","visibleOnTemplate":1,"visibleOnTerminal":0,"filePath":"","ApplicationId":"763d0d8f-a0fd-4fa6-97e3-e44028305ba3","ParamCategoryId":"6ba90b8c-0340-44fe-82c4-4bc32e51a316","categoryName":"AMP Cloud","updatedAt":"2021-11-18T06:18:22.000Z","createdAt":"2021-11-18T06:18:22.000Z"}]}}`)
		}

	})

	tp := TemplateParams{}
	opt := ParamsOpt{
		CategoryId: "value1",
	}
	err := c.TemplatesService.GetParams(context.Background(), "814", opt, &tp)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 1

	if tp.Count != want {
		t.Errorf("Templates count = %v, want %v", tp.Count, want)
	}
}

func TestTemplatesUpdateParamsMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		want := "value2"
		if !templateRe.MatchString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			testMethod(t, r, http.MethodPost)
			if r.PostFormValue("param2") != want {
				t.Errorf("incorrect form value got %v, want %v", r.PostFormValue("param2"), want)
			}
			fmt.Fprint(w, `{"success":true,"message":"Successfully updated 3 parameter(s) and propagated the changes to 2 terminals.","failed":["string"],"updated":["string"]}`)
		}
	})

	wantErr := errors.New("required templateID is missing")

	updated := []string{}
	failed := []string{}
	params := map[string]string{}
	files := map[string]string{}

	err := c.TemplatesService.UpdateParams(context.Background(), "", params, files, &updated, &failed)
	if err == nil {
		t.Errorf("Error is nil, want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("Error got %v, want %v", err, wantErr)
	}

	params["param1"] = "value1"
	params["param2"] = "value2"
	params["param3"] = "value3"
	files["TODO"] = "./TODO"
	err = c.TemplatesService.UpdateParams(context.Background(), "814", params, files, &updated, &failed)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}
	if updated[0] != "string" {
		t.Errorf("updated is incorrect got %v, want \"string\"", updated[0])
	}
	if failed[0] != "string" {
		t.Errorf("failed is incorrect got %v, want \"string\"", failed[0])
	}
}
