package parser

import (
	"encoding/json"
)

func SortJson(s []byte) ([]byte, error) {
	var foundObj map[string]interface{}
	unmarshalErr := json.Unmarshal(s, &foundObj)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	// json.Marshal sorts keys
	sorted, marshalErr := json.Marshal(foundObj)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return sorted, nil
}
