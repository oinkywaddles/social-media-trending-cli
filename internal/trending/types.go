package trending

import "time"

type RequestPolicy struct {
	Upstream         string `json:"upstream"`
	LocalMinInterval string `json:"local_min_interval"`
	CacheTTL         string `json:"cache_ttl"`
}

type Item struct {
	Rank        int            `json:"rank"`
	Title       string         `json:"title"`
	Score       string         `json:"score,omitempty"`
	ScoreValue  *int64         `json:"score_value,omitempty"`
	Label       string         `json:"label,omitempty"`
	Meta        string         `json:"meta,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	URL         string         `json:"url,omitempty"`
	Cover       string         `json:"cover,omitempty"`
	PublishedAt string         `json:"published_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
	Raw         map[string]any `json:"raw,omitempty"`
}

type Result struct {
	Provider        string        `json:"provider"`
	Platform        Platform      `json:"platform"`
	DisplayName     string        `json:"display_name"`
	Endpoint        string        `json:"endpoint"`
	FetchedAt       time.Time     `json:"fetched_at"`
	Cached          bool          `json:"cached"`
	TotalCount      int           `json:"total_count"`
	UpstreamMessage string        `json:"upstream_message,omitempty"`
	RequestPolicy   RequestPolicy `json:"request_policy"`
	Items           []Item        `json:"items"`
}
