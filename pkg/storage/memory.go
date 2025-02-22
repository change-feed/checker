package storage

import (
	"errors"
	"sync"
)

// MemoryStorage is a mock StorageBackend for testing
type MemoryStorage struct {
	data map[string]map[string][]byte
	mu   sync.RWMutex
}

// NewMemoryStorage initializes a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]map[string][]byte),
	}
}

func (m *MemoryStorage) List(resourceKey string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if snapshots, exists := m.data[resourceKey]; exists {
		ids := make([]string, 0, len(snapshots))
		for id := range snapshots {
			ids = append(ids, id)
		}
		return ids, nil
	}

	return []string{}, nil
}

func (m *MemoryStorage) Load(resourceKey, id string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if snapshots, exists := m.data[resourceKey]; exists {
		if snapshot, found := snapshots[id]; found {
			return snapshot, nil
		}
	}

	return nil, errors.New("snapshot not found")
}

func (m *MemoryStorage) Save(resourceKey, id string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[resourceKey]; !exists {
		m.data[resourceKey] = make(map[string][]byte)
	}
	m.data[resourceKey][id] = data

	return nil
}

func (m *MemoryStorage) Delete(resourceKey, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if snapshots, exists := m.data[resourceKey]; exists {
		delete(snapshots, id)
		if len(snapshots) == 0 {
			delete(m.data, resourceKey)
		}
		return nil
	}

	return errors.New("snapshot not found")
}
