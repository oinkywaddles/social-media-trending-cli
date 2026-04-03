package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"social-media-trending-cli/internal/trending"
)

type jsonEnvelope struct {
	GeneratedAt time.Time    `json:"generated_at"`
	Results     []jsonResult `json:"results"`
}

type jsonResult struct {
	Platform      trending.Platform `json:"platform"`
	FetchedAt     time.Time         `json:"fetched_at"`
	TotalCount    int               `json:"total_count"`
	ReturnedCount int               `json:"returned_count"`
	Items         []jsonItem        `json:"items"`
}

type jsonItem struct {
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Score       string `json:"score,omitempty"`
	Label       string `json:"label,omitempty"`
	Meta        string `json:"meta,omitempty"`
	URL         string `json:"url,omitempty"`
	Summary     string `json:"summary,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
}

type jsonRawEnvelope struct {
	GeneratedAt time.Time         `json:"generated_at"`
	Results     []trending.Result `json:"results"`
}

type jsonDetailEnvelope struct {
	GeneratedAt time.Time         `json:"generated_at"`
	Platform    trending.Platform `json:"platform"`
	FetchedAt   time.Time         `json:"fetched_at"`
	TotalCount  int               `json:"total_count"`
	Item        jsonDetailItem    `json:"item"`
}

type jsonDetailItem struct {
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Score       string `json:"score,omitempty"`
	Label       string `json:"label,omitempty"`
	Meta        string `json:"meta,omitempty"`
	URL         string `json:"url,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Cover       string `json:"cover,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type jsonDetailRawEnvelope struct {
	GeneratedAt     time.Time              `json:"generated_at"`
	Provider        string                 `json:"provider"`
	Platform        trending.Platform      `json:"platform"`
	DisplayName     string                 `json:"display_name"`
	Endpoint        string                 `json:"endpoint"`
	FetchedAt       time.Time              `json:"fetched_at"`
	Cached          bool                   `json:"cached"`
	TotalCount      int                    `json:"total_count"`
	UpstreamMessage string                 `json:"upstream_message,omitempty"`
	RequestPolicy   trending.RequestPolicy `json:"request_policy"`
	Item            trending.Item          `json:"item"`
}

func PrintJSON(w io.Writer, results []trending.Result) {
	out := jsonEnvelope{
		GeneratedAt: time.Now().UTC(),
		Results:     make([]jsonResult, 0, len(results)),
	}

	for _, result := range results {
		row := jsonResult{
			Platform:      result.Platform,
			FetchedAt:     result.FetchedAt,
			TotalCount:    result.TotalCount,
			ReturnedCount: len(result.Items),
			Items:         make([]jsonItem, 0, len(result.Items)),
		}
		for _, item := range result.Items {
			row.Items = append(row.Items, jsonItem{
				Rank:        item.Rank,
				Title:       item.Title,
				Score:       item.Score,
				Label:       item.Label,
				Meta:        item.Meta,
				URL:         item.URL,
				Summary:     item.Summary,
				PublishedAt: item.PublishedAt,
			})
		}
		out.Results = append(out.Results, row)
	}

	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Fprintln(w, string(data))
}

func PrintJSONRaw(w io.Writer, results []trending.Result) {
	rawData, _ := json.MarshalIndent(jsonRawEnvelope{
		GeneratedAt: time.Now().UTC(),
		Results:     results,
	}, "", "  ")
	fmt.Fprintln(w, string(rawData))
}

func PrintDetailJSON(w io.Writer, result trending.Result, item trending.Item) {
	data, _ := json.MarshalIndent(jsonDetailEnvelope{
		GeneratedAt: time.Now().UTC(),
		Platform:    result.Platform,
		FetchedAt:   result.FetchedAt,
		TotalCount:  result.TotalCount,
		Item: jsonDetailItem{
			Rank:        item.Rank,
			Title:       item.Title,
			Score:       item.Score,
			Label:       item.Label,
			Meta:        item.Meta,
			URL:         item.URL,
			Summary:     item.Summary,
			Cover:       item.Cover,
			PublishedAt: item.PublishedAt,
			UpdatedAt:   item.UpdatedAt,
		},
	}, "", "  ")
	fmt.Fprintln(w, string(data))
}

func PrintDetailJSONRaw(w io.Writer, result trending.Result, item trending.Item) {
	data, _ := json.MarshalIndent(jsonDetailRawEnvelope{
		GeneratedAt:     time.Now().UTC(),
		Provider:        result.Provider,
		Platform:        result.Platform,
		DisplayName:     result.DisplayName,
		Endpoint:        result.Endpoint,
		FetchedAt:       result.FetchedAt,
		Cached:          result.Cached,
		TotalCount:      result.TotalCount,
		UpstreamMessage: result.UpstreamMessage,
		RequestPolicy:   result.RequestPolicy,
		Item:            item,
	}, "", "  ")
	fmt.Fprintln(w, string(data))
}
