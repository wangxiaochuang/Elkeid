package safemap

import "sync"

type SafeMap struct {
	name    string
	dataMap map[string]map[string]interface{}
	mu      sync.RWMutex
}

func NewSafeMap(name string) *SafeMap {
	sm := &SafeMap{
		name:    name,
		dataMap: make(map[string]map[string]interface{}),
	}
	return sm
}

func (sm *SafeMap) Name() string {
	return sm.name
}

func (sm *SafeMap) HSet(key, subKey string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.dataMap[key]; !ok {
		sm.dataMap[key] = make(map[string]interface{})
	}
	sm.dataMap[key][subKey] = value
}

func (sm *SafeMap) HGet(key, subKey string) interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if _, ok := sm.dataMap[key]; ok {
		if _, ok = sm.dataMap[key][subKey]; ok {
			return sm.dataMap[key][subKey]
		}
	}
	return nil
}

func (sm *SafeMap) HDel(key, subKey string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.dataMap[key]; ok {
		if _, ok = sm.dataMap[key][subKey]; ok {
			delete(sm.dataMap[key], subKey)
			if len(sm.dataMap[key]) == 0 {
				delete(sm.dataMap, key)
			}
		}
	}
}

func (sm *SafeMap) Del(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.dataMap, key)
}

func (sm *SafeMap) Get(key string) map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	r := make(map[string]interface{})
	if _, ok := sm.dataMap[key]; ok {
		for k, v := range sm.dataMap[key] {
			r[k] = v
		}
		return r
	}
	return nil
}

func (sm *SafeMap) Keys() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	r := make([]string, 0)
	for k, _ := range sm.dataMap {
		r = append(r, k)
	}
	return r
}

func (sm *SafeMap) HKeys(key string) []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	r := make([]string, 0)
	if _, ok := sm.dataMap[key]; ok {
		for k, _ := range sm.dataMap[key] {
			r = append(r, k)
		}
	}
	return r
}

func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.dataMap)
}

func (sm *SafeMap) HLen(key string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if _, ok := sm.dataMap[key]; ok {
		return len(sm.dataMap[key])
	} else {
		return 0
	}
}
