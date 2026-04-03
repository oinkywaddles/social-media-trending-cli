package cmd

import (
	"fmt"
	"os"

	"social-media-trending-cli/internal/output"
	"social-media-trending-cli/internal/trending"

	"github.com/spf13/cobra"
)

var (
	flagLimit int
)

var listCmd = &cobra.Command{
	Use:   "list [platform ...|all]",
	Short: "Fetch trending lists for one or more platforms",
	Args:  cobra.ArbitraryArgs,
	Run:   runList,
}

func init() {
	listCmd.Flags().IntVar(&flagLimit, "limit", 0, "Limit items per platform")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	validateOutputFlags()

	platforms, err := trending.ResolvePlatforms(args)
	if err != nil {
		exitWithError(err.Error())
	}

	provider := buildProvider()
	store := buildSnapshotStore()

	results := make([]trending.Result, 0, len(platforms))
	for _, platform := range platforms {
		result, err := provider.Fetch(cmd.Context(), platform)
		if err != nil {
			exitWithError(fmt.Sprintf("%s: %v", platform, err))
		}
		if _, err := store.Save(result); err != nil {
			warn("Warning: failed to save snapshot: " + err.Error())
		}

		snapshotResult := result
		if flagLimit > 0 && len(result.Items) > flagLimit {
			result.Items = result.Items[:flagLimit]
		}
		_ = snapshotResult
		results = append(results, result)
	}

	if jsonRawOutput {
		output.PrintJSONRaw(os.Stdout, results)
		return
	}
	if jsonOutput {
		output.PrintJSON(os.Stdout, results)
		return
	}
	output.PrintText(os.Stdout, results)
}
