package guid

import (
	"sync"

	. "github.com/nece099/base/logger"
)

type GUIDManager struct {
	storage *sync.Map
}

var guidManager = &GUIDManager{
	storage: &sync.Map{},
}

func GetGUIDManager() *GUIDManager {
	return guidManager
}

func (g *GUIDManager) Save(obj interface{}) int64 {
	newid := GenerateUID()
	g.storage.Store(newid, obj)
	return newid
}

func (g *GUIDManager) Load(id int64) interface{} {
	obj, ok := g.storage.Load(id)
	if !ok {
		Log.Errorf("id=[%v] does not exist in GUID storage", id)
		return nil
	}

	return obj
}

func (g *GUIDManager) Remove(id int64) {
	g.storage.Delete(id)
}
