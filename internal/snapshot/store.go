package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"social-media-trending-cli/internal/trending"
)

const DefaultTTL = time.Hour

type Entry struct {
	StoredAt time.Time       `json:"stored_at"`
	Result   trending.Result `json:"result"`
}

type Store struct {
	dir string
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func DefaultDir() string {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		return filepath.Join(os.TempDir(), "social-media-trending-cli", "snapshots")
	}
	return filepath.Join(base, "social-media-trending-cli", "snapshots")
}

func (s *Store) LoadFresh(platform trending.Platform, ttl time.Duration) (Entry, bool, error) {
	if s == nil || ttl <= 0 {
		return Entry{}, false, nil
	}

	data, err := os.ReadFile(s.pathForPlatform(platform))
	if err != nil {
		if os.IsNotExist(err) {
			return Entry{}, false, nil
		}
		return Entry{}, false, err
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, false, err
	}
	if time.Since(entry.StoredAt) > ttl {
		return Entry{}, false, nil
	}
	return entry, true, nil
}

func (s *Store) Save(result trending.Result) (Entry, error) {
	if s == nil {
		return Entry{StoredAt: time.Now().UTC(), Result: result}, nil
	}

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return Entry{}, err
	}

	entry := Entry{
		StoredAt: time.Now().UTC(),
		Result:   result,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, err
	}
	if err := os.WriteFile(s.pathForPlatform(result.Platform), data, 0o644); err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (s *Store) pathForPlatform(platform trending.Platform) string {
	return filepath.Join(s.dir, platform.String()+".json")
}
