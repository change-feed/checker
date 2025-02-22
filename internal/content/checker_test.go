package content

import (
	"testing"

	"github.com/change-feed/checker/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotChecker(t *testing.T) {
	mockBackend := storage.NewMemoryStorage()

	checker, err := NewContentChecker(mockBackend, "https://example.com")
	assert.NoError(t, err)

	changed, err := checker.Check()
	assert.NoError(t, err)
	assert.True(t, changed)

	// Run again; should not detect a change
	changed, err = checker.Check()
	assert.NoError(t, err)
	assert.False(t, changed)
}
