package amp360

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingRoundTripper struct {
	Wrapped http.RoundTripper
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	fmt.Printf("Resquest: | %v | %v | \n", req.Method, req.URL)

	start := time.Now()
	res, err = l.Wrapped.RoundTrip(req)
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		fmt.Printf("Response: | %v | %v |\n", res.Status, time.Since(start))
	}

	return res, err
}
