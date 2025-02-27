package service

import "sync"

// Cache stores already processed image paths
var Cache = struct {
	sync.RWMutex
	data map[string]map[string]string
}{data: make(map[string]map[string]string)}

// GetCachedResult returns the cached result if available
func GetCachedResult(filename string) (map[string]string, bool) {
	Cache.RLock()
	defer Cache.RUnlock()
	result, exists := Cache.data[filename]
	return result, exists
}

// CacheResult saves processed image paths
func CacheResult(filename string, paths map[string]string) {
	Cache.Lock()
	defer Cache.Unlock()
	Cache.data[filename] = paths
}
