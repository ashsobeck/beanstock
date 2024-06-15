package tests

import (
	"beanstock/internal/diff"
	"beanstock/internal/server"
	"beanstock/internal/types"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestHandler(t *testing.T) {
	s := &server.Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func Test_HashDiff__NoDifference(t *testing.T) {
	jsonStr := `{
		"Name": "noknow",
		"Admin": true,
		"Hobbies": ["IT","Travel"],
		"Address": {
		    "PostalCode": 1111,
		    "Country": "Japan"
		},
		"Null": null,
		"Age": 2
	    }`

	newJsonStr := `{
		"Name": "noknow",
		"Null": null,
		"Hobbies": ["IT","Travel"],
		"Address": {
		    "PostalCode": 1111,
		    "Country": "Japan"
		},
		"Admin": true,
		"Age": 2
	    }`

	var obj map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &obj)

	lastJson, _ := json.Marshal(obj)
	lastHash, _ := sha256.New().Write(lastJson)

	site := types.Website{
		Id:           uuid.NewString(),
		Url:          "",
		ShopProvider: "",
		Json:         map[string]interface{}{},
		LastHash:     lastHash,
	}

	res, err := diff.HashDiff(site, []byte(newJsonStr))
	if res != false || err != nil {
		t.Fatalf("should result in no difference; %v err: %v", res, err)
	}
}

func Test_HashDiff__Difference(t *testing.T) {
	jsonStr := `{
		"Name": "noknow",
		"Admin": false,
		"Hobbies": ["IT","Travel"],
		"Address": {
		    "PostalCode": 1111,
		    "Country": "Japan"
		},
		"Null": null,
		"Age": 2
	    }`

	newJsonStr := `{
		"Name": "noknow",
		"Null": null,
		"Hobbies": ["IT","Travel"],
		"Address": {
		    "PostalCode": 1111,
		    "Country": "Japan"
		},
		"Admin": true,
		"Age": 2
	    }`

	var obj map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &obj)

	lastJson, _ := json.Marshal(obj)
	lastHash, _ := sha256.New().Write(lastJson)

	site := types.Website{
		Id:           uuid.NewString(),
		Url:          "",
		ShopProvider: "",
		Json:         map[string]interface{}{},
		LastHash:     lastHash,
	}

	res, err := diff.HashDiff(site, []byte(newJsonStr))
	if res != true || err != nil {
		t.Fatalf("should result in difference; %v err: %v", res, err)
	}
}
