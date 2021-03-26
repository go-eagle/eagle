package cache

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var storeObj *store

type store struct {
	sync.RWMutex
	cron *cron.Cron
	data map[int]interface{}
}

func NewSyncStore() *store {
	return &store{
		data: make(map[int]interface{}, 0),
		cron: cron.New(),
	}
}

func (s *store) syncStore(dataFn func() map[int]interface{}) {
	s.Lock()
	s.data = dataFn()
	s.Unlock()
}

func (s *store) Get(id int) interface{} {
	s.RLock()
	defer s.RUnlock()
	resp, _ := s.data[id]
	return resp
}
