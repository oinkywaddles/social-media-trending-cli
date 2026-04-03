package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
)

var (
	version = "dev"

	jsonOutput      bool
	jsonRawOutput   bool
	noCache         bool
	timeoutFlag     time.Duration
	cacheTTLFlag    time.Duration
	minIntervalFlag time.Duration
)

func getVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

var rootCmd = &cobra.Command{
	Use:     "social-media-trending-cli",
	Short:   "Fetch social media trending lists across major Chinese platforms",
	Version: getVersion(),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output trimmed JSON for downstream consumption")
	rootCmd.PersistentFlags().BoolVar(&jsonRawOutput, "json-raw", false, "Output full JSON including provider metadata and raw upstream payload")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "Disable local cache")
	rootCmd.PersistentFlags().DurationVar(&timeoutFlag, "timeout", 15*time.Second, "HTTP timeout per request")
	rootCmd.PersistentFlags().DurationVar(&cacheTTLFlag, "cache-ttl", 2*time.Minute, "Local cache TTL")
	rootCmd.PersistentFlags().DurationVar(&minIntervalFlag, "min-interval", 1200*time.Millisecond, "Minimum interval between upstream requests")
}

func exitWithError(msg string) {
	if usesJSONOutput() {
		data, _ := json.Marshal(map[string]string{"error": msg})
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		fmt.Fprintln(os.Stderr, "Error: "+msg)
	}
	os.Exit(1)
}

func usesJSONOutput() bool {
	return jsonOutput || jsonRawOutput
}

func validateOutputFlags() {
	if jsonOutput && jsonRawOutput {
		exitWithError("use either --json or --json-raw, not both")
	}
}
