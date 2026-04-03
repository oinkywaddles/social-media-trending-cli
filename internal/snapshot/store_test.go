package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"social-media-trending-cli/internal/trending"
)

func TestStoreSaveAndLoadFresh(t *testing.T) {
	store := NewStore(t.TempDir())
	result := trending.Result{
		Platform:    trending.PlatformZhihu,
		DisplayName: "Zhihu",
		FetchedAt:   time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
		Items: []trending.Item{
			{Rank: 1, Title: "topic"},
		},
	}

	saved, err := store.Save(result)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if saved.StoredAt.IsZero() {
		t.Fatal("expected stored_at to be set")
	}

	loaded, fresh, err := store.LoadFresh(trending.PlatformZhihu, time.Hour)
	if err != nil {
		t.Fatalf("LoadFresh returned error: %v", err)
	}
	if !fresh {
		t.Fatal("expected snapshot to be fresh")
	}
	if loaded.Result.Platform != trending.PlatformZhihu {
		t.Fatalf("expected platform zhihu, got %q", loaded.Result.Platform)
	}
	if len(loaded.Result.Items) != 1 || loaded.Result.Items[0].Title != "topic" {
		t.Fatalf("unexpected loaded items: %+v", loaded.Result.Items)
	}
}

func TestStoreLoadFreshExpiresByStoredAt(t *testing.T) {
	store := NewStore(t.TempDir())
	entry := Entry{
		StoredAt: time.Now().Add(-2 * time.Hour).UTC(),
		Result: trending.Result{
			Platform: trending.PlatformWeibo,
		},
	}

	dataPath := filepath.Join(store.dir, "weibo.json")
	data, err := jsonMarshal(entry)
	if err != nil {
		t.Fatalf("marshal entry: %v", err)
	}
	if err := osWriteFile(dataPath, data); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}

	_, fresh, err := store.LoadFresh(trending.PlatformWeibo, time.Hour)
	if err != nil {
		t.Fatalf("LoadFresh returned error: %v", err)
	}
	if fresh {
		t.Fatal("expected expired snapshot to be treated as stale")
	}
}

var (
	jsonMarshal = func(v any) ([]byte, error) {
		return json.Marshal(v)
	}
	osWriteFile = func(path string, data []byte) error {
		return os.WriteFile(path, data, 0o644)
	}
)
