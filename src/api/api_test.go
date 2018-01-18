package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/placer14/ob-crawler/api"
	"gopkg.in/jarcoal/httpmock.v1"
)

func newTestClient() *api.Client {
	return &api.Client{
		Client:     &http.Client{},
		Host:       "foo",
		Port:       123,
		AuthCookie: "bar",
	}
}

func equalHashSlices(s1, s2 []string) bool {
	if s1 == nil && s2 == nil {
		return true
	}
	if s1 == nil || s2 == nil {
		return false
	}
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func setup() {
	httpmock.Activate()
}

func teardown() {
	httpmock.DeactivateAndReset()
}

func registerCookieValidatingResponder(url, expectedCookie string, response *http.Response, t *testing.T) {
	var cookieMatched = false
	httpmock.RegisterResponder(http.MethodGet, url,
		func(req *http.Request) (*http.Response, error) {
			for _, cookieValue := range req.Header["Cookie"] {
				if expectedCookie == cookieValue {
					cookieMatched = true
				}
			}
			if cookieMatched == false {
				t.Errorf("Expected cookie to be included, but was not\n\tActual headers: %+v", req.Header)
			}
			return response, nil
		},
	)
}

func TestGetPeersIsRequestedAndProcessedSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	var (
		c             = newTestClient()
		testPath      = fmt.Sprintf("http://%s:%d/ob/peers", c.Host, c.Port)
		expectedPeers = []string{
			"QmRDcEDK9gSViAevCHiE6ghkaBCU7rTuQj4BDpmCzRvRYg",
			"QmUZRGLhcKXF1JyuaHgKm23LvqcoMYwtb9jmh8CkP4og3K",
		}
	)

	response, err := httpmock.NewJsonResponse(http.StatusOK, expectedPeers)
	if err != nil {
		t.Error(err)
	}
	registerCookieValidatingResponder(testPath, c.AuthCookie, response, t)

	actualPeers, err := c.GetPeers()
	if err != nil {
		t.Error(err)
	}
	if equalHashSlices(expectedPeers, actualPeers) != true {
		t.Errorf("Expected responses to match\n\tExpected: %+v\n\tActual: %+v\n", expectedPeers, actualPeers)
	}
}

func TestGetClosestPeersIsRequestedAndProcessedSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	var (
		c             = newTestClient()
		targetPeer    = "peer"
		testPath      = fmt.Sprintf("http://%s:%d/ob/closestpeers/%s", c.Host, c.Port, targetPeer)
		expectedPeers = []string{
			"QmRDcEDK9gSViAevCHiE6ghkaBCU7rTuQj4BDpmCzRvRYg",
			"QmUZRGLhcKXF1JyuaHgKm23LvqcoMYwtb9jmh8CkP4og3K",
		}
	)

	response, err := httpmock.NewJsonResponse(http.StatusOK, expectedPeers)
	if err != nil {
		t.Error(err)
	}
	registerCookieValidatingResponder(testPath, c.AuthCookie, response, t)

	actualPeers, err := c.GetClosestPeers(targetPeer)
	if err != nil {
		t.Error(err)
	}
	if equalHashSlices(expectedPeers, actualPeers) != true {
		t.Errorf("Expected responses to match\n\tExpected: %+v\n\tActual: %+v\n", expectedPeers, actualPeers)
	}
}

func TestGetListingsCountIsRequestedAndProcessedSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	var (
		c          = newTestClient()
		targetPeer = "peer"
		testPath   = fmt.Sprintf("http://%s:%d/ob/profile/%s", c.Host, c.Port, targetPeer)
		// a list of anonymous things
		expectedCount = 5
		stats         = &api.StatisticsStub{ListingCount: expectedCount}
		profile       = api.ProfileStub{PeerID: targetPeer, Stats: stats}
	)

	response, err := httpmock.NewJsonResponse(http.StatusOK, profile)
	if err != nil {
		t.Error(err)
	}
	registerCookieValidatingResponder(testPath, c.AuthCookie, response, t)

	count, err := c.GetListingsCount(targetPeer)
	if err != nil {
		t.Error(err)
	}
	if count != expectedCount {
		t.Error("Expected response to be properly counted, but was not")
	}
}

func TestGetListingsCountReturnsZeroWhenListingsNotFound(t *testing.T) {
	setup()
	defer teardown()

	var (
		c          = newTestClient()
		targetPeer = "peer"
		testPath   = fmt.Sprintf("http://%s:%d/ob/listings/%s", c.Host, c.Port, targetPeer)
	)

	response, err := httpmock.NewJsonResponse(http.StatusNotFound, []interface{}{})
	if err != nil {
		t.Error(err)
	}
	registerCookieValidatingResponder(testPath, c.AuthCookie, response, t)

	count, err := c.GetListingsCount(targetPeer)
	if err == nil {
		t.Error("Expected an error, where there was actually none")
	}
	if count != 0 {
		t.Error("Expected response to be zero, but was not")
	}
}
