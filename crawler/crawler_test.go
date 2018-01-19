package crawler

import (
	"strings"
	"sync"
	"testing"

	"github.com/placer14/ob-crawler/api"
)

func prettyCalls(calls []*api.CallEntry, joiner string) string {
	var callStrings []string
	for _, call := range calls {
		callStrings = append(callStrings, call.String())
	}
	return strings.Join(callStrings, joiner)
}

func TestExpectedCallsAreMadeToTheAPI(t *testing.T) {
	crawler := &Crawler{
		api: &api.FakeClient{
			MethodResponses: map[string]interface{}{
				"getpeers":        []string{"peer1", "peer2", "peer3"},
				"getclosestpeers": []string{"peer4", "peer5", "peer3"},
				"getlistingscount": map[string]int{
					"peer1": 1,
					"peer2": 1,
					"peer3": 1,
					"peer4": 1,
					"peer5": 1,
				},
			},
		},
		cacheMutex:     &sync.Mutex{},
		workerPoolSize: 1,
		workersActive:  &sync.WaitGroup{},
	}
	crawler.Execute()

	// Expected each peer to GetClosestPeers and GetListingsCount
	// plus one call to GetPeers to seed the work queue
	// Despite each call to GetCloestPeers appending its list of
	// repeated peers, this asserts it will only ever be called once.
	var (
		expectedCallEntries = []*api.CallEntry{
			&api.CallEntry{"getpeers", ""},
			&api.CallEntry{"getclosestpeers", "peer1"},
			&api.CallEntry{"getclosestpeers", "peer2"},
			&api.CallEntry{"getclosestpeers", "peer3"},
			&api.CallEntry{"getclosestpeers", "peer4"},
			&api.CallEntry{"getclosestpeers", "peer5"},
			&api.CallEntry{"getlistingscount", "peer1"},
			&api.CallEntry{"getlistingscount", "peer2"},
			&api.CallEntry{"getlistingscount", "peer3"},
			&api.CallEntry{"getlistingscount", "peer4"},
			&api.CallEntry{"getlistingscount", "peer5"},
		}
		actualCallEntries = (crawler.api.(*api.FakeClient)).CallRecord
	)

	for _, expectedCall := range expectedCallEntries {
		var expectationPassed = false
		for _, recordedCall := range actualCallEntries {
			if expectedCall.EqualTo(recordedCall) {
				expectationPassed = true
			}
		}
		if expectationPassed == false {
			t.Errorf("Missing call for '%s' with arg '%s'\n", expectedCall.Method, expectedCall.Arg)
		}
	}

	if len(expectedCallEntries) != len(actualCallEntries) {
		t.Errorf("Number of call entries do not match")
		t.Logf("Expected:\n%+v \n\nActual:\n%+v\n\n", prettyCalls(expectedCallEntries, "\n"), prettyCalls(actualCallEntries, "\n"))
	}
}

func TestListingCountIsAccurate(t *testing.T) {
	getListingCounts := map[string]int{
		"peer1": 1,
		"peer2": 2,
		"peer3": 3,
		"peer4": 4,
		"peer5": 5,
	}

	crawler := &Crawler{
		api: &api.FakeClient{
			MethodResponses: map[string]interface{}{
				"getpeers":         []string{"peer1", "peer2", "peer3"},
				"getclosestpeers":  []string{"peer4", "peer5", "peer3"},
				"getlistingscount": getListingCounts,
			},
		},
		cacheMutex:     &sync.Mutex{},
		workerPoolSize: 1,
		workersActive:  &sync.WaitGroup{},
	}
	crawler.Execute()

	expectedTotal := 0
	for _, count := range getListingCounts {
		expectedTotal += count
	}
	actualTotal := crawler.ListingCount()
	if expectedTotal != actualTotal {
		t.Errorf("Expected to equal %d, but was %d\n", expectedTotal, actualTotal)
	}
}

func TestNodesVisitedIsAccurate(t *testing.T) {
	getListingCounts := map[string]int{
		"peer1": 1,
		"peer2": 1,
		"peer3": 1,
		"peer4": 1,
		"peer5": 1,
	}

	crawler := &Crawler{
		api: &api.FakeClient{
			MethodResponses: map[string]interface{}{
				"getpeers":         []string{"peer1", "peer2", "peer3"},
				"getclosestpeers":  []string{"peer4", "peer5", "peer3"},
				"getlistingscount": getListingCounts,
			},
		},
		cacheMutex:           &sync.Mutex{},
		workerPoolSize:       1,
		workersActive:        &sync.WaitGroup{},
		maximumVisitsAllowed: 2,
	}
	crawler.Execute()

	expectedTotal := crawler.maximumVisitsAllowed
	actualTotal := crawler.NodesVisited()
	if expectedTotal != actualTotal {
		t.Errorf("Expected to equal %d, but was %d\n", expectedTotal, actualTotal)
	}
}
