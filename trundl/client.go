package trundl

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://trundl.awstrp.net/"
)

var token = os.Getenv("TRUNDL_SECRET")

type roundTripperFunc func(r *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type Client struct {
	token   string
	baseURL *url.URL
	client  *http.Client
}

func NewClient(httpClient *http.Client, opts ...func(*Client)) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	u, _ := url.Parse(defaultBaseURL)
	c := &Client{baseURL: u, token: token, client: httpClient}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithRTLogger(rt http.RoundTripper) func(*Client) {
	return func(c *Client) {

		if rt == nil {
			rt = http.DefaultTransport
		}

		c.client.Transport = roundTripperFunc(func(req *http.Request) (resp *http.Response, err error) {
			log.Printf("[%p]: %s %d %s\n", req, req.Method, resp.StatusCode, req.URL.Path)
			return rt.RoundTrip(req)
		})
	}
}

func (c *Client) NewRequest(ctx context.Context, endpoint string, vs url.Values) (*http.Request, error) {
	rel, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = vs.Encode()
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}
	return resp, err
}

func (c *Client) Records(ctx context.Context, ids []int) (*Records, error) {
	beg, end := timeInterval()
	vs := make(url.Values)
	vs.Set("begin_date", beg)
	vs.Set("end_date", end)
	vs.Set("price_type", "CLOSE")
	vs.Set("security_unique_identifier", intsToString(ids))
	vs.Set("token", c.token) // Likely is no token; likely need to confirm auth on request header
	req, err := c.NewRequest(ctx, "api/data/IDW_SECURITY_PRICE", vs)
	if err != nil {
		return nil, err
	}
	var rx Records
	_, err = c.Do(req, &rx)
	return &rx, err
}

func timeInterval() (string, string) {
	t := time.Now()
	layoutISO := "2006-01-02"
	beg := t.AddDate(0, 0, -8).Format(layoutISO)
	end := t.AddDate(0, 0, -1).Format(layoutISO)
	return beg, end
}

func intsToString(ix []int) string {
	b := make([]string, len(ix))
	for i, v := range ix {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, ",")
}
