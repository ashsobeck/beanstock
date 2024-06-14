package diff

import (
	"beanstock/internal/database"
	"beanstock/internal/types"
	"crypto/sha256"
	"io"
	"net/http"
)

type Differ interface {
	HashDiff(site types.Website) (bool, error)
}

type differ struct{}

func (d *differ) HashDiff(site types.Website) (bool, error) {
	res, err := http.Get(site.Url)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	json, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return false, readErr
	}

	newHash, hashErr := sha256.New().Write(json)
	if hashErr != nil {
		return false, hashErr
	}

	if newHash != site.LastHash {
		return true, nil
	}

	//TODO: store new hash

	return false, nil
}
