// Memorydb provides a key value store.
//nolint:golint,stylecheck // TODO
package memorydb

import (
	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

/*

In memory Key/Value store
- Load from flat file
- Store to flat file when exiting
- Store to flat file periodically as backup

Methods:
- Create
- Read/Get
- Update
- Delete

*/

type MemoryStore struct {
	*cache.Cache
}

func New() *MemoryStore {
	return &MemoryStore{
		cache.New(cache.NoExpiration, 0),
	}
}

func LoadDBFrom(dbLocation string) error {
	logrus.Println("Loaded memory DB from file: ", dbLocation)
	return nil
}
