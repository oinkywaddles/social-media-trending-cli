package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type fileEntry struct {
	FetchedAt time.Time       `json:"fetched_at"`
	Body      json.RawMessage `json:"body"`
}

type FileCache struct {
	dir string
}

func NewFileCache(dir string) *FileCache {
	return &FileCache{dir: dir}
}

func DefaultDir() string {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		return filepath.Join(os.TempDir(), "social-media-trending-cli")
	}
	return filepath.Join(base, "social-media-trending-cli")
}

func (c *FileCache) ReadFresh(key string, ttl time.Duration) ([]byte, time.Time, bool, error) {
	if c == nil || ttl <= 0 {
		return nil, time.Time{}, false, nil
	}

	path := c.pathForKey(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, time.Time{}, false, nil
		}
		return nil, time.Time{}, false, err
	}

	var entry fileEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, time.Time{}, false, err
	}
	if time.Since(entry.FetchedAt) > ttl {
		return nil, time.Time{}, false, nil
	}
	return entry.Body, entry.FetchedAt, true, nil
}

func (c *FileCache) Write(key string, body []byte, fetchedAt time.Time) error {
	if c == nil {
		return nil
	}

	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}

	entry := fileEntry{
		FetchedAt: fetchedAt.UTC(),
		Body:      append([]byte(nil), body...),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(c.pathForKey(key), data, 0o644)
}

func (c *FileCache) pathForKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(c.dir, hex.EncodeToString(sum[:])+".json")
}
