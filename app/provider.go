package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	IP   string `json:"Ip"`
	Port int    `json:"Port"`
}

type Client struct {
	host   *url.URL
	client http.Client
}

func NewClient(rawURL string, timeout time.Duration) (*Client, error) {
	host, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse credit-fs client URL: %w`, err)
	}
	client := http.Client{
		Timeout: timeout,
	}

	return &Client{
		client: client,
		host:   host,
	}, nil
}

func (c Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	URL, err := c.host.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, URL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c Client) doRequest(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) GetListOfProxies() (Response, error) {
	req, err := c.newRequest(http.MethodGet, "/api/proxy", nil)
	if err != nil {
		return Response{}, err
	}

	q := req.URL.Query()

	q.Add("format", "json")
	q.Add("limit", "1")
	q.Add("uptime", "50")
	q.Add("type", "socks5")
	q.Add("ping", "500")
	
	req.URL.RawQuery = q.Encode()

	res, err := c.doRequest(req)
	if err != nil {
		return Response{}, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Response{}, errors.New("failed to get")
	}

	var dataString string
	if err = json.NewDecoder(res.Body).Decode(&dataString); err != nil {
		return Response{}, err
	}

	var response []Response
	err = json.Unmarshal([]byte(dataString), &response)
	if err != nil {
		return Response{}, err
	}

	if len(response) < 1 {
		return Response{}, errors.New("not found")
	}

	return response[0], nil
}
