package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/form"
)

var encoder *form.Encoder

func init() {
	encoder = form.NewEncoder()
}

// Client used to make api calls to the AT API server
type Client struct {
	username   string
	apiKey     string
	sandbox    bool
	httpClient *http.Client
}

// New creates a new Client app.
// Sandbox indicates whether the requests will be targeted to the sandbox AT endpoint
func New(apiKey, username string, sandbox bool) *Client {
	cli := &http.Client{}
	return &Client{apiKey: apiKey, username: username, sandbox: sandbox, httpClient: cli}
}

// SetHTTPClient overrides the default http client
func (c *Client) SetHTTPClient(cli *http.Client) {
	c.httpClient = cli
}

// Endpoint represent the different services exposed by AT
type Endpoint string

// URL constructs a full valid url with the host
func (e Endpoint) URL(sandbox bool) string {
	var host string
	switch sandbox {
	case true:
		host = SandboxHost
	case false:
		host = LiveHost
	}

	return fmt.Sprintf("https://%s%s", host, e)
}

const (
	// V1EndpointMessaging for sending SMS
	V1EndpointMessaging Endpoint = "/version1/messaging"
	// LiveHost for live endpoint host
	LiveHost = "api.africastalking.com"
	// SandboxHost for sandbox endpont host
	SandboxHost = "api.sandbox.africastalking.com"
)

// UsernameSetter allows the caller to set the username from AT
type UsernameSetter interface {
	SetUsername(username string)
}

// Do makes the API call and returns a reader containing the AT response body
func (c *Client) Do(req interface{}, endpoint Endpoint) (rep io.Reader, err error) {
	return c.do(req, endpoint)
}

func (c *Client) do(reqBody interface{}, endpoint Endpoint) (repReader io.Reader, err error) {
	// if reqBody implements UsernameSetter, set it
	if setter, ok := reqBody.(UsernameSetter); ok {
		setter.SetUsername(c.username)
	}

	values, err := encoder.Encode(reqBody)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(nil)
	_, err = body.WriteString(values.Encode())
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.URL(c.sandbox), body)
	if err != nil {
		panic(err)
	}
	// Setup necessary headers i.e apiKey, Content-Type and Accept
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", c.apiKey)

	rep, err := c.httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	return rep.Body, nil
}
