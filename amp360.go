package amp360

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.1"
	defaultBase    = "https://api.amp360.amobilepayment.com/v1/"
	devBase        = "https://dev.api.amp360.amobilepayment.com/v1/"
	defaultUA      = "go-amp360-client/" + libraryVersion
)

var (
	ErrUnauthorized = errors.New("api: unauthorized access")
	ErrNotFound     = errors.New("api: not found")
)

func NewClient(defaultBaseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if defaultBaseURL == "" {
		defaultBaseURL = defaultBase
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		BaseURL:   baseURL,
		UserAgent: defaultUA,
		client:    httpClient,
	}
	c.TemplatesService = &TemplatesService{client: c}
	c.CompaniesService = &CompaniesService{client: c}
	c.ModelsService = &ModelsService{client: c}
	c.TerminalsService = &TerminalsService{client: c}
	return c
}

type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	BaseURL   *url.URL
	UserAgent string
	apiKey    string

	TemplatesService *TemplatesService
	CompaniesService *CompaniesService
	ModelsService    *ModelsService
	TerminalsService *TerminalsService
}

type service struct {
	client *Client
}

func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// Client returns the http.Client used by this AMP360 client.
func (c *Client) Client() *http.Client {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	clientCopy := *c.client
	return &clientCopy
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	return c.NewRequestCtx(context.Background(), method, path, body)
}

func (c *Client) NewMultiPartRequestCtx(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body.(io.Reader))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Authorization", c.apiKey)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) NewRequestCtx(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Authorization", c.apiKey)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	return resp, err
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	vs, err := query.Values(opt)
	if err != nil {
		return s, err
	}
	u.RawQuery = vs.Encode()
	return u.String(), nil
}
