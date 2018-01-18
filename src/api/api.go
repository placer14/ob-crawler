package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	*http.Client
	Host       string
	Port       int
	AuthCookie string
}

func New() *Client {
	return &Client{
		Client: &http.Client{},
	}
}

func (c *Client) GetPeers() ([]string, error) {
	var (
		err         error
		peers       []string
		req         *http.Request
		requestPath string
	)

	requestPath = fmt.Sprintf("http://%s:%d/ob/%s", c.Host, c.Port, "peers")
	req, err = http.NewRequest(http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %s\n", err)
	}
	req.Header.Add("cookie", c.AuthCookie)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %s\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: status %s\n", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&peers); err != nil {
		return nil, fmt.Errorf("decoding: %s\n", err)
	}
	return peers, nil
}
