package cmd

import (
	"os"

	"social-media-trending-cli/internal/cache"
	"social-media-trending-cli/internal/httpx"
	"social-media-trending-cli/internal/snapshot"
	"social-media-trending-cli/internal/trending"
)

func buildProvider() trending.Provider {
	var fileCache *cache.FileCache
	cacheTTL := cacheTTLFlag
	if !noCache {
		fileCache = cache.NewFileCache(cache.DefaultDir())
	} else {
		cacheTTL = 0
	}

	client := httpx.NewClient(timeoutFlag, minIntervalFlag, fileCache)
	return trending.NewSixtySProvider(client, cacheTTL, minIntervalFlag)
}

func buildSnapshotStore() *snapshot.Store {
	return snapshot.NewStore(snapshot.DefaultDir())
}

func warn(msg string) {
	if !usesJSONOutput() {
		_, _ = os.Stderr.WriteString(msg + "\n")
	}
}
