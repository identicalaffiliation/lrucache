// LEAST RECENTLY USED CACHE IMPLEMENTATION.
//
// Cache interface is defined in interface.go.
// Thread-safe implementation using sync.RWMutex.
// No external dependencies.
//
// Usage:
//
//	cache := NewCache(capacity)
//
//	cache.Set(key, value)
//
//	value, ok := cache.Get(key)
//	If the key is not found, value = -1 and ok = false.
//
// Thread safety:
//   - Set and Get use an exclusive lock because they modify the cache.
//   - Len uses a read lock because it only reads the cache state.

package lrucache
