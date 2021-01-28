package controllerstest

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/vsokoltsov/beeg/tests"
)

// TestMain rewrite base test for events controller module
func TestMain(m *testing.M) {
	tests.TestMain(m)
}

// TestSuccessEventCreationRoute tests success creation of event
func TestSuccessEventCreationRoute(t *testing.T) {
	key := "1-test"
	tests.DeleteFromRedis(key)
	defer tests.DeleteFromRedis(key)

	var data = map[string]interface{}{
		"id":    1,
		"label": "test",
	}
	var resposneMap = make(map[string]json.RawMessage)
	d, _ := json.Marshal(data)
	response := tests.MakeRequest("POST", "/", strings.NewReader(string(d)))
	json.Unmarshal(response.Body.Bytes(), &resposneMap)
	if string(resposneMap["viewed"]) != "1" {
		t.Error("Viewed value for these id and label does not match")
	}
}

// TestSuccessEventCreationRoute tests failed creation of event
func TestFailedEventCreationRoute(t *testing.T) {
	key := "1-test"
	tests.DeleteFromRedis(key)
	defer tests.DeleteFromRedis(key)

	var data = map[string]interface{}{}
	var resposneMap = make(map[string]json.RawMessage)
	d, _ := json.Marshal(data)
	response := tests.MakeRequest("POST", "/", strings.NewReader(string(d)))
	json.Unmarshal(response.Body.Bytes(), &resposneMap)
	if resposneMap["errors"] == nil {
		t.Error("Viewed value for these id and label does not match")
	}
}
