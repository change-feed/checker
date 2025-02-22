package content

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/change-feed/checker/internal/util"
	"github.com/change-feed/checker/pkg/storage"
)

// SnapshotChecker encapsulates snapshot processing logic
type SnapshotChecker struct {
	backend        storage.StorageBackend
	resourceKey    string
	compressedData []byte
	snapshots      []string
}

// NewContentChecker initializes a new checker for a given URL and backend
func NewContentChecker(backend storage.StorageBackend, url string) (*SnapshotChecker, error) {
	slog.Debug("Preparing to check URL for changes...", "url", url)

	resourceKey, compressedData, err := prepare(url)
	if err != nil {
		return nil, err
	}

	return &SnapshotChecker{
		backend:        backend,
		resourceKey:    resourceKey,
		compressedData: compressedData,
	}, nil
}

// Check determines if a change has occurred and updates snapshots
func (c *SnapshotChecker) Check() (bool, error) {
	snapshots, err := c.backend.List(c.resourceKey)
	if err != nil {
		return false, fmt.Errorf("failed to refresh snapshots: %w", err)
	}
	c.snapshots = snapshots

	if changed, err := c.hasChanged(); err != nil || !changed {
		return changed, err
	}

	if err := c.pruneSnapshots(); err != nil {
		return false, err
	}

	if err := c.storeSnapshot(); err != nil {
		return false, err
	}

	slog.Debug("Change detected", "url", c.resourceKey)
	return true, nil
}

// hasChanged compares the latest snapshot to the new one
func (c *SnapshotChecker) hasChanged() (bool, error) {
	if len(c.snapshots) == 0 {
		return true, nil // No snapshots exist yet, so it's a new change
	}

	latestSnapshotID := c.snapshots[len(c.snapshots)-1]
	latestSnapshot, err := c.backend.Load(c.resourceKey, latestSnapshotID)
	if err != nil {
		return false, fmt.Errorf("error loading snapshot: %w", err)
	}

	return util.Hash(latestSnapshot) != util.Hash(c.compressedData), nil
}

// pruneSnapshots deletes the oldest snapshot if storage exceeds the limit
func (c *SnapshotChecker) pruneSnapshots() error {
	if len(c.snapshots) > util.MaxSnapshots {
		oldestSnapshotID := c.snapshots[0]
		if err := c.backend.Delete(c.resourceKey, oldestSnapshotID); err != nil {
			return fmt.Errorf("error deleting snapshot: %w", err)
		}
	}

	return nil
}

// storeSnapshot saves the new snapshot
func (c *SnapshotChecker) storeSnapshot() error {
	id := time.Now().Format("20060102-150405")
	if err := c.backend.Save(c.resourceKey, id, c.compressedData); err != nil {
		return fmt.Errorf("error saving snapshot: %w", err)
	}

	return nil
}
