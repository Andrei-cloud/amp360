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
	terminalRe      = regexp.MustCompile(`^\/terminals\/params\/(\d+)`)
	terminalBulkRe  = regexp.MustCompile(`^\/terminals\/params\/bulk\/(\d+)`)
	terminalReQuery = regexp.MustCompile(`^\/terminals\/params\/(\d+)\?categoryId=([a-z0-9\-]+)\&?`)
)

func TestTerminalGetParamsMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !terminalRe.MatchString(r.URL.Path) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"success":true,"message":"Successfully found the terminal parameters.","data":{"categories":[{"name":"TERMINAL","id":"c2f7c244-ebd7-4ce3-bcf6-adffd2e4ec90"}],"count":2,"rows":[{"id":3104508,"type":"STRING","tag":"ACQS._1.ACQINFO.MERCHANTID","name":"ACQS._1.ACQINFO.MERCHANTID","hint":"","validator":"","value":"400081203","defaultValue":"000000000","visibleOnTemplate":1,"visibleOnTerminal":1,"filePath":null,"ApplicationId":"c250f201-4d0b-42d3-aeb0-3c8804e4684a","ParamCategoryId":"c2f7c244-ebd7-4ce3-bcf6-adffd2e4ec90","categoryName":"TERMINAL"},{"id":3104509,"type":"STRING","tag":"ACQS._1.ACQINFO.MERCHANTNAME","name":"ACQS._1.ACQINFO.MERCHANTNAME","hint":"","validator":"","value":"MRA RESTAURANT BAKERY & SWEETS","defaultValue":"NONE","visibleOnTemplate":1,"visibleOnTerminal":1,"filePath":null,"ApplicationId":"c250f201-4d0b-42d3-aeb0-3c8804e4684a","ParamCategoryId":"c2f7c244-ebd7-4ce3-bcf6-adffd2e4ec90","categoryName":"TERMINAL"}]}}`)
		}
	})

	tp := TerminalParams{}
	wantErr := errors.New("required terminalID is missing")
	err := c.TerminalsService.GetParams(context.Background(), 0, nil, &tp)
	if err == nil {
		t.Errorf("Error is nil, want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("Error got %v, want %v", err, wantErr)
	}

	want := 1
	err = c.TerminalsService.GetParams(context.Background(), 814, nil, &tp)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}
	if len(tp.Categories) != want {
		t.Errorf("Categories count got %v, wnat %v", len(tp.Categories), want)
	}
	if tp.Count != 2 {
		t.Errorf("Parameters count got %v, wnat %v", len(tp.Categories), 2)
	}
	if len(tp.Rows) != 2 {
		t.Errorf("Parameters actual count got %v, wnat %v", len(tp.Rows), 2)
	}
}

func TestTerminalGetParamsQueryMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !terminalReQuery.MatchString(r.URL.String()) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"message":"bad request","data":{}}`)
			t.Errorf("Bad URL got %v", r.URL.Path)
		} else {
			testMethod(t, r, http.MethodGet)
			values := r.URL.Query().Get("categoryId")
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
	err := c.TerminalsService.GetParams(context.Background(), 814, opt, &tp)
	if err != nil {
		t.Errorf("Error occured = %v", err)
	}

	want := 1

	if tp.Count != want {
		t.Errorf("Terminal parameters count = %v, want %v", tp.Count, want)
	}
}

func TestTerminalsUpdateParamsMock(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	c.client.Transport = LoggingRoundTripper{http.DefaultTransport}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		want := "value2"
		if !terminalBulkRe.MatchString(r.URL.Path) {
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

	wantErr := errors.New("required terminalID is missing")

	updated := []string{}
	failed := []string{}
	params := map[string]string{}
	files := map[string]string{}

	err := c.TerminalsService.UpdateParams(context.Background(), 0, params, files, &updated, &failed)
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
	err = c.TerminalsService.UpdateParams(context.Background(), 814, params, files, &updated, &failed)
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
