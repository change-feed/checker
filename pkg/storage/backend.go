package storage

type StorageBackend interface {
	List(resourceKey string) ([]string, error)      // Retrieves all snapshot IDs, ordered by ascending timestamp
	Load(resourceKey, id string) ([]byte, error)    // Retrieves a specific snapshot
	Save(resourceKey, id string, data []byte) error // Stores a new snapshot
	Delete(resourceKey, id string) error            // Deletes a specific snapshot
}
