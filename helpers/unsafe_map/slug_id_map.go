package unsafe_map

import "sync"

var SlugIDMap = slugIdMap{
	m:   make(map[string]int32, 4096),
	rwm: sync.RWMutex{},
}

type slugIdMap struct {
	m   map[string]int32
	rwm sync.RWMutex
}

func (sim *slugIdMap) StoreSlug(slug string, id int32) {
	sim.rwm.Lock()
	sim.m[slug] = id
	sim.rwm.Unlock()
}

func (sim *slugIdMap) GetId(slug string) int32 {
	return sim.m[slug]
}

func (sim *slugIdMap) Clear() {
	sim.m = make(map[string]int32, 4096)
}
