package content

import (
	"fmt"
	"io"
	"net/http"

	"github.com/change-feed/checker/internal/util"
)

// prepare fetches and compresses the content for the given URL
func prepare(url string) (string, []byte, error) {
	data, err := fetch(url)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	compressedData, err := util.Compress(data)
	if err != nil {
		return "", nil, fmt.Errorf("failed to compress data: %w", err)
	}

	resourceKey := util.SanitizeString(url)

	return resourceKey, compressedData, nil
}

// fetch retrieves the content from the given URL
func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
