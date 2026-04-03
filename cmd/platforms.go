package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"social-media-trending-cli/internal/trending"

	"github.com/spf13/cobra"
)

var platformsCmd = &cobra.Command{
	Use:   "platforms",
	Short: "List supported platforms and upstream endpoints",
	Run:   runPlatforms,
}

func init() {
	rootCmd.AddCommand(platformsCmd)
}

func runPlatforms(cmd *cobra.Command, args []string) {
	validateOutputFlags()

	infos := trending.SupportedPlatformInfos()
	if usesJSONOutput() {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Fprintln(os.Stdout, string(data))
		return
	}

	fmt.Fprintln(os.Stdout, "| Platform | Upstream Endpoint | Aliases |")
	fmt.Fprintln(os.Stdout, "|----------|-------------------|---------|")
	for _, info := range infos {
		fmt.Fprintf(
			os.Stdout,
			"| %s | %s | %s |\n",
			info.DisplayName,
			info.Endpoint,
			strings.Join(info.Aliases, ", "),
		)
	}
}
