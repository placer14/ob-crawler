package api

import "fmt"

// CallEntry captures the details of call made to a FakeClient method
type CallEntry struct {
	Method string
	Arg    string
}

// EqualTo tests the equality of two CallEntry instances
func (c *CallEntry) EqualTo(c2 *CallEntry) bool {
	return c.Method == c2.Method && c.Arg == c2.Arg
}

// String returns a pretty, easy-to-read collection of letters for
// human consumption
func (c *CallEntry) String() string { return fmt.Sprintf("%s(%s)", c.Method, c.Arg) }

// FakeClient is a testable mock of the OpenBazaar API. The mock response
// can be assigned on MethodResponses with the key as the lowercased method
// name and the value as the desired return value. Each stubbed method
// is responsible for casting to the method's appropriate return type.
type FakeClient struct {
	CallRecord      []*CallEntry
	MethodResponses map[string]interface{}
}

// Below are stubbed methods for FakeClient to satisfy api.OpenBazaarlikeAPI
// interface and mirror the external behavior defined in api.Client

func (f *FakeClient) HostPort() string { return "fake.api.client:123" }

func (f *FakeClient) GetPeers() ([]string, error) {
	f.CallRecord = append(f.CallRecord, &CallEntry{"getpeers", ""})
	response := f.MethodResponses["getpeers"].([]string)
	return response, nil
}

func (f *FakeClient) GetClosestPeers(hash string) ([]string, error) {
	f.CallRecord = append(f.CallRecord, &CallEntry{"getclosestpeers", hash})
	response := f.MethodResponses["getclosestpeers"].([]string)
	return response, nil
}

func (f *FakeClient) GetListingsCount(hash string) (int, error) {
	f.CallRecord = append(f.CallRecord, &CallEntry{"getlistingscount", hash})
	responses := f.MethodResponses["getlistingscount"].(map[string]int)
	return responses[hash], nil
}
