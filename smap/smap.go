package smap

import "sync"

// Map is an implementation of a synchronized map[string]string
type Map struct {
	entries map[string]string
	mutex   sync.RWMutex
}

// New returns a new Map
func New() *Map {
	return &Map{
		entries: make(map[string]string),
	}
}

// Get retrieves a key's value and whether or not it exists
func (m *Map) Get(key string) (string, bool) {
	m.lock(false)
	defer m.unlock(false)
	return m.get(key)
}

func (m *Map) get(key string) (string, bool) {
	value, exists := m.entries[key]
	if !exists {
		return "", exists
	}
	return value, exists
}

// Delete will remove a value from the map and return whether or not it existed
func (m *Map) Delete(key string) bool {
	m.lock(true)
	defer m.unlock(true)
	return m.delete(key)
}

func (m *Map) delete(key string) bool {
	contains := m.contains(key)
	delete(m.entries, key)
	return contains
}

// Put adds a value to the map and returns if it was actually an update
func (m *Map) Put(key, value string) bool {
	m.lock(true)
	defer m.unlock(true)
	return m.put(key, value)
}

func (m *Map) put(key, value string) bool {
	updated := m.contains(key)
	m.entries[key] = value
	return updated
}

// Replace will change the value if it exists
func (m *Map) Replace(key, value string) bool {
	m.lock(true)
	defer m.unlock(true)
	return m.replace(key, value)
}

func (m *Map) replace(key, value string) bool {
	if m.contains(key) {
		return m.put(key, value)
	}
	return false
}

// Alter will apply fn to the key's value, if it exists
func (m *Map) Alter(key string, fn func(string) string) bool {
	m.lock(true)
	defer m.unlock(true)
	return m.alter(key, fn)
}

func (m *Map) alter(key string, fn func(string) string) bool {
	v, contains := m.get(key)
	if !contains {
		return contains
	}
	m.put(key, fn(v))
	return contains
}

// Contains -- whether or not that map has this key
func (m *Map) Contains(key string) bool {
	m.lock(false)
	defer m.unlock(false)
	return m.contains(key)
}

func (m *Map) contains(key string) bool {
	_, contains := m.entries[key]
	return contains
}

// ContainsValue -- whether or not the map has the value
func (m *Map) ContainsValue(search string) bool {
	m.lock(false)
	defer m.unlock(false)

	for _, value := range m.entries {
		if search == value {
			return true
		}
	}
	return false
}

// Size returns the number of entries in the map
func (m *Map) Size() int {
	m.lock(false)
	size := len(m.entries)
	m.unlock(false)
	return size
}

// IsEmpty is true if Size()
func (m *Map) IsEmpty() bool {
	return m.Size() == 0
}

func (m *Map) forEach(fn func(string, string) bool) {
	for key, value := range m.entries {
		if fn(key, value) {
			return
		}
	}
}

// Merge will combine m2 into the smap, and return a slice of keys that were updated
func (m *Map) Merge(m2 *Map) []string {
	m.lock(true)
	m2.lock(true)
	defer m.unlock(true)
	defer m2.unlock(true)

	var additions []string
	for key, value := range m2.entries {
		if m.contains(key) {
			additions = append(additions, key)
		}
		m.put(key, value)
	}
	return additions
}

// Transform will change every key's value using the function
func (m *Map) Transform(fn func(string) string) {
	m.lock(true)
	defer m.unlock(true)

	m.forEach(func(key, value string) bool {
		m.put(key, fn(value))
		return false
	})
}

func (m *Map) lock(write bool) {
	if write {
		m.mutex.Lock()
		return
	}

	m.mutex.RLock()
}

func (m *Map) unlock(write bool) {
	if write {
		m.mutex.Unlock()
		return
	}

	m.mutex.RUnlock()
}
