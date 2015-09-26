package safemap

import (
	"sync"

	"github.com/cagnosolutions/safemap/util"
)

var SHARD_COUNT int

type SafeMap []*Shard

type Shard struct {
	items map[string]interface{}
	sync.RWMutex
}

func NewSafeMap(shardCount int) *SafeMap {
	if shardCount == 0 || shardCount%2 != 0 {
		shardCount = 16
	}
	SHARD_COUNT = shardCount
	m := make(SafeMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &Shard{
			items: make(map[string]interface{}),
		}
	}
	return &m
}

func (m *SafeMap) GetShard(key string) *Shard {
	bucket := util.Sum32([]byte(key)) % uint32(SHARD_COUNT)
	//fmt.Printf("key: %q, bucket: %d\n", key, bucket) // <= Remove this at some point
	return (*m)[bucket]
}

func (m *SafeMap) Set(key string, val interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = val
	shard.Unlock()
}

func (m *SafeMap) Get(key string) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m *SafeMap) Del(key string) {
	if shard := m.GetShard(key); shard != nil {
		shard.Lock()
		delete(shard.items, key)
		shard.Unlock()
	}
}

type EntrySet struct {
	Key string
	Val interface{}
}

func (m *SafeMap) Iter() <-chan EntrySet {
	ch := make(chan EntrySet)
	go func() {
		for _, shard := range *m {
			shard.RLock()
			for key, val := range shard.items {
				ch <- EntrySet{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}
