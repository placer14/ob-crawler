package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenBazaarlikeAPI interface {
	HostPort() string
	GetPeers() ([]string, error)
	GetClosestPeers(string) ([]string, error)
	GetListingsCount(string) (int, error)
}

// Client is a configurable abstraction around the OpenBazaar API which handles
// adding the authentication header to endpoints of interest
type Client struct {
	*http.Client
	Host       string
	Port       int
	AuthCookie string
}

// New returns a pointer to a new Client
func New() *Client {
	return &Client{
		Client: &http.Client{},
	}
}

func (c *Client) HostPort() string { return fmt.Sprintf("%s:%d", c.Host, c.Port) }

func (c *Client) doRequest(endpoint string) (*http.Response, error) {
	requestPath := fmt.Sprintf("http://%s:%d/ob/%s", c.Host, c.Port, endpoint)
	req, err := http.NewRequest(http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %s\n", err)
	}
	req.Header.Add("Cookie", c.AuthCookie)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status not ok: %s\n", resp.Status)
	}
	return resp, nil
}

// GetClosestPeers retrieves the list of peers of the node referenced by
// the hash originHash
func (c *Client) GetClosestPeers(originHash string) ([]string, error) {
	resp, err := c.doRequest(fmt.Sprintf("closestpeers/%s", originHash))
	if err != nil {
		return nil, fmt.Errorf("requesting closest peers: %s\n", err)
	}

	peers := make([]string, 0)
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = decoder.Decode(&peers); err != nil {
		return nil, fmt.Errorf("decoding closest peers: %s\n", err)
	}
	return peers, nil
}

// GetPeers retrieves the list of immediate peers of the API node
func (c *Client) GetPeers() ([]string, error) {
	resp, err := c.doRequest("peers")
	if err != nil {
		return nil, fmt.Errorf("requesting peers: %s\n", err)
	}

	peers := make([]string, 0)
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = decoder.Decode(&peers); err != nil {
		return nil, fmt.Errorf("decoding peers: %s\n", err)
	}
	return peers, nil
}

// GetListingsCount retrieves the list of contracts offered by the node
// referenced by the hash originHash
func (c *Client) GetListingsCount(originHash string) (int, error) {
	resp, err := c.doRequest(fmt.Sprintf("profile/%s", originHash))
	if err != nil {
		return 0, fmt.Errorf("requesting profile: %s\n", err)
	}

	profile := &ProfileStub{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = decoder.Decode(profile); err != nil {
		return 0, fmt.Errorf("decoding profile: %s\n", err)
	}
	return profile.Stats.ListingCount, nil
}

type StatisticsStub struct {
	ListingCount int `json:"listingCount"`
}
type ProfileStub struct {
	PeerID string          `json:"peerID"`
	Stats  *StatisticsStub `json:"stats"`
}
