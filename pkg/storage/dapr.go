package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dapr/go-sdk/client"
)

type DaprStorageBackend struct {
	daprClient  client.Client
	stateStore  string
}

// NewDaprStorage initializes the Dapr-backed state storage
func NewDaprStorage() (StorageBackend, error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Dapr client: %w", err)
	}
	return &DaprStorageBackend{
		daprClient: daprClient,
		stateStore: "snapshot-store",
	}, nil
}

// List retrieves all snapshot IDs for a given resourceKey
func (s *DaprStorageBackend) List(resourceKey string) ([]string, error) {
	query := fmt.Sprintf(`{"filter": {"key": {"prefix": "%s/"}}}`, resourceKey)
	meta := map[string]string{}

	queryResp, err := s.daprClient.QueryStateAlpha1(context.TODO(), s.stateStore, query, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshots: %w", err)
	}

	var snapshotIDs []string
	for _, item := range queryResp.Results {
		snapshotID := strings.TrimPrefix(item.Key, resourceKey+"/")
		snapshotIDs = append(snapshotIDs, snapshotID)
	}

	// Ensure snapshots are ordered by timestamp (if snapshot IDs are time-based)
	sort.Strings(snapshotIDs)
	return snapshotIDs, nil
}

// Load retrieves a specific snapshot
func (s *DaprStorageBackend) Load(resourceKey, id string) ([]byte, error) {
	snapshotKey := fmt.Sprintf("%s/%s", resourceKey, id)
	meta := map[string]string{}

	item, err := s.daprClient.GetState(context.TODO(), s.stateStore, snapshotKey, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve snapshot: %w", err)
	}
	if item.Value == nil {
		return nil, fmt.Errorf("snapshot not found: %s", snapshotKey)
	}

	return item.Value, nil
}

// Save stores a new snapshot
func (s *DaprStorageBackend) Save(resourceKey, id string, data []byte) error {
	snapshotKey := fmt.Sprintf("%s/%s", resourceKey, id)
	meta := map[string]string{}

	err := s.daprClient.SaveState(context.TODO(), s.stateStore, snapshotKey, data, meta)
	if err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	return nil
}

// Delete removes a specific snapshot
func (s *DaprStorageBackend) Delete(resourceKey, id string) error {
	snapshotKey := fmt.Sprintf("%s/%s", resourceKey, id)
	meta := map[string]string{}

	err := s.daprClient.DeleteState(context.TODO(), s.stateStore, snapshotKey, meta)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	return nil
}
