package output

import (
	"fmt"
	"io"
	"strings"

	"social-media-trending-cli/internal/trending"
)

func PrintText(w io.Writer, results []trending.Result) {
	for i, result := range results {
		if i > 0 {
			fmt.Fprintln(w)
		}

		fmt.Fprintf(w, "%s Trending\n", result.DisplayName)
		fmt.Fprintf(
			w,
			"Source: %s (%s) | showing %d of %d | fetched at %s | cached=%t\n",
			result.Provider,
			result.Endpoint,
			len(result.Items),
			result.TotalCount,
			result.FetchedAt.Format("2006-01-02 15:04:05 MST"),
			result.Cached,
		)
		fmt.Fprintf(
			w,
			"Request policy: %s | local throttle %s | cache TTL %s\n\n",
			result.RequestPolicy.Upstream,
			result.RequestPolicy.LocalMinInterval,
			result.RequestPolicy.CacheTTL,
		)

		if len(result.Items) == 0 {
			fmt.Fprintln(w, "No trending items found.")
			continue
		}

		fmt.Fprintln(w, "| # | Title | Score | Meta |")
		fmt.Fprintln(w, "|---:|-------|-------|------|")
		for _, item := range result.Items {
			fmt.Fprintf(
				w,
				"| %d | %s | %s | %s |\n",
				item.Rank,
				markdownCell(item.Title, 48),
				markdownCell(fallback(item.Score, "-"), 18),
				markdownCell(buildMeta(item), 36),
			)
		}
	}
}

func PrintDetailText(w io.Writer, result trending.Result, item trending.Item) {
	fmt.Fprintf(w, "%s #%d\n", result.DisplayName, item.Rank)
	fmt.Fprintf(
		w,
		"Fetched at: %s | Source: %s (%s)\n",
		result.FetchedAt.Format("2006-01-02 15:04:05 MST"),
		result.Provider,
		result.Endpoint,
	)
	fmt.Fprintf(w, "Title: %s\n", item.Title)

	if line := joinNonEmpty(
		prefixed("Score", item.Score),
		prefixed("Label", item.Label),
		prefixed("Meta", item.Meta),
	); line != "" {
		fmt.Fprintln(w, line)
	}
	if item.PublishedAt != "" {
		fmt.Fprintf(w, "Published: %s\n", item.PublishedAt)
	}
	if item.UpdatedAt != "" {
		fmt.Fprintf(w, "Updated: %s\n", item.UpdatedAt)
	}
	if item.URL != "" {
		fmt.Fprintf(w, "URL: %s\n", item.URL)
	}
	if item.Cover != "" {
		fmt.Fprintf(w, "Cover: %s\n", item.Cover)
	}
	if item.Summary != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Summary:")
		fmt.Fprintln(w, item.Summary)
	}
}

func buildMeta(item trending.Item) string {
	parts := make([]string, 0, 3)
	if item.Label != "" {
		parts = append(parts, item.Label)
	}
	if item.Meta != "" {
		parts = append(parts, item.Meta)
	}
	if len(parts) == 0 && item.URL != "" {
		parts = append(parts, item.URL)
	}
	return strings.Join(parts, " | ")
}

func markdownCell(value string, limit int) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return truncateRunes(value, limit)
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if limit <= 3 {
		return string(runes[:limit])
	}
	return string(runes[:limit-3]) + "..."
}

func fallback(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func prefixed(label, value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return label + ": " + strings.TrimSpace(value)
}

func joinNonEmpty(values ...string) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			parts = append(parts, strings.TrimSpace(value))
		}
	}
	return strings.Join(parts, " | ")
}
