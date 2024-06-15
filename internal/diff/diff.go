package diff

import (
	"beanstock/internal/types"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	// "io"
	// "net/http"
)

type Differ interface {
	HashDiff(site types.Website) (bool, error)
}

type differ struct{}

func HashDiff(site types.Website, body []byte) (bool, error) {
	// res, err := http.Get(site.Url)
	// if err != nil {
	// 	return false, err
	// }
	//
	// defer res.Body.Close()
	//
	// body, readErr := io.ReadAll(res.Body)
	// if readErr != nil {
	// 	return false, readErr
	// }

	var foundObj map[string]interface{}
	unmarshalErr := json.Unmarshal(body, &foundObj)
	if unmarshalErr != nil {
		return false, unmarshalErr
	}

	// json.Marshal sorts keys
	sorted, marshalErr := json.Marshal(foundObj)
	if marshalErr != nil {
		return false, marshalErr
	}

	fmt.Println(string(sorted))

	newHash, hashErr := sha256.New().Write(sorted)
	if hashErr != nil {
		return false, hashErr
	}

	fmt.Println(site.LastHash, newHash)

	if newHash != site.LastHash {
		return true, nil
	}

	//TODO: store new hash

	return false, nil
}
