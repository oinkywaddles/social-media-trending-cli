package cmd

import (
	"fmt"
	"os"
	"strconv"

	"social-media-trending-cli/internal/output"
	"social-media-trending-cli/internal/snapshot"
	"social-media-trending-cli/internal/trending"

	"github.com/spf13/cobra"
)

var detailCmd = &cobra.Command{
	Use:   "detail <platform> <rank>",
	Short: "Show detail for a trending item from the latest platform snapshot",
	Args:  cobra.ExactArgs(2),
	Run:   runDetail,
}

func init() {
	rootCmd.AddCommand(detailCmd)
}

func runDetail(cmd *cobra.Command, args []string) {
	validateOutputFlags()

	platform, err := trending.ParsePlatform(args[0])
	if err != nil {
		exitWithError(err.Error())
	}

	rank, err := strconv.Atoi(args[1])
	if err != nil || rank <= 0 {
		exitWithError("rank must be a positive integer")
	}

	provider := buildProvider()
	store := buildSnapshotStore()

	result, err := loadOrRefreshSnapshot(cmd, provider, store, platform)
	if err != nil {
		exitWithError(fmt.Sprintf("%s: %v", platform, err))
	}

	item, ok := findItemByRank(result, rank)
	if !ok {
		exitWithError(fmt.Sprintf("%s rank %d not found; available range is 1-%d", platform, rank, len(result.Items)))
	}

	if jsonRawOutput {
		output.PrintDetailJSONRaw(os.Stdout, result, item)
		return
	}
	if jsonOutput {
		output.PrintDetailJSON(os.Stdout, result, item)
		return
	}
	output.PrintDetailText(os.Stdout, result, item)
}

func loadOrRefreshSnapshot(cmd *cobra.Command, provider trending.Provider, store *snapshot.Store, platform trending.Platform) (trending.Result, error) {
	entry, fresh, err := store.LoadFresh(platform, snapshot.DefaultTTL)
	if err != nil {
		warn("Warning: failed to read snapshot, refreshing live data: " + err.Error())
	}
	if err == nil && fresh {
		return entry.Result, nil
	}

	result, err := provider.Fetch(cmd.Context(), platform)
	if err != nil {
		return trending.Result{}, err
	}
	if _, err := store.Save(result); err != nil {
		warn("Warning: failed to save snapshot: " + err.Error())
	}
	return result, nil
}

func findItemByRank(result trending.Result, rank int) (trending.Item, bool) {
	for _, item := range result.Items {
		if item.Rank == rank {
			return item, true
		}
	}
	if rank >= 1 && rank <= len(result.Items) {
		return result.Items[rank-1], true
	}
	return trending.Item{}, false
}
