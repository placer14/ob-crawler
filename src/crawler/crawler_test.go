package crawler

import (
	"strings"
	"sync"
	"testing"

	"github.com/placer14/ob-crawler/api"
)

func newTestCrawler() *Crawler {
	return &Crawler{
		api: &api.FakeClient{
			MethodResponses: map[string]interface{}{
				"getpeers":         []string{"peer1", "peer2", "peer3"},
				"getclosestpeers":  []string{"peer4", "peer5", "peer3"},
				"getlistingscount": 1,
			},
		},
		cacheMutex:     &sync.Mutex{},
		workerPoolSize: 1,
		workersActive:  &sync.WaitGroup{},
	}
}

func prettyCalls(calls []*api.CallEntry, joiner string) string {
	var callStrings []string
	for _, call := range calls {
		callStrings = append(callStrings, call.String())
	}
	return strings.Join(callStrings, joiner)
}

func TestExpectedCallsAreMadeToTheAPI(t *testing.T) {
	crawler := newTestCrawler()
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
