package store

import (
	"sync"

	"github.com/cagnosolutions/safemap"
)

type SafeMapStore struct {
	SafeMaps map[string]*safemap.SafeMap
	sync.RWMutex
}

func NewSafeMapStore(shardCount int) *SafeMapStore {
	if shardCount == 0 || shardCount%2 != 0 {
		shardCount = 16
	}
	safemap.SHARD_COUNT = shardCount
	return &SafeMapStore{
		SafeMaps: make(map[string]*safemap.SafeMap),
	}
}

func (sms *SafeMapStore) Set(key, fld string, val interface{}) {
	sm, ok := sms.GetSafeMap(key)
	if !ok {
		sms.Lock()
		sms.SafeMaps[key] = safemap.NewSafeMap(safemap.SHARD_COUNT)
		sms.Unlock()
	}
	sm.Set(fld, val)
}

func (sms *SafeMapStore) Get(key, fld string) (interface{}, bool) {
	if sm, ok := sms.GetSafeMap(key); ok {
		return sm.Get(fld)
	}
	return nil, false
}

func (sms *SafeMapStore) Del(key, fld string) {
	if sm, ok := sms.GetSafeMap(key); ok {
		sm.Del(fld)
	}
}

func (sms *SafeMapStore) AddStore(key string) {
	if _, ok := sms.GetSafeMap(key); !ok {
		sms.Lock()
		sms.SafeMaps[key] = safemap.NewSafeMap(safemap.SHARD_COUNT)
		sms.Unlock()
	}
}

func (sms *SafeMapStore) GetSafeMap(key string) (*safemap.SafeMap, bool) {
	sms.RLock()
	sm, ok := sms.SafeMaps[key]
	sms.RUnlock()
	return sm, ok
}

func (sms *SafeMapStore) DelStore(key string) {
	if _, ok := sms.GetSafeMap(key); ok {
		sms.Lock()
		delete(sms.SafeMaps, key)
		sms.Unlock()
	}
}
