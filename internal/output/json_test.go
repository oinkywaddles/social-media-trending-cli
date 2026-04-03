package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"social-media-trending-cli/internal/trending"
)

func TestPrintJSONTrimsRawFields(t *testing.T) {
	results := []trending.Result{
		{
			Platform:   trending.PlatformZhihu,
			FetchedAt:  time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
			TotalCount: 30,
			Items: []trending.Item{
				{
					Rank:        1,
					Title:       "topic",
					Score:       "525 万热度",
					Meta:        "236 answers | 553 followers",
					URL:         "https://example.com",
					PublishedAt: "2026/04/03 13:29:58",
					Summary:     "very long summary",
					Raw:         map[string]any{"detail": "raw detail"},
				},
			},
		},
	}

	var buf bytes.Buffer
	PrintJSON(&buf, results)
	out := buf.String()

	if !strings.Contains(out, `"platform": "zhihu"`) {
		t.Fatalf("expected platform in lite json, got %s", out)
	}
	if strings.Contains(out, `"raw"`) {
		t.Fatalf("did not expect raw payload in lite json, got %s", out)
	}
	if !strings.Contains(out, `"summary": "very long summary"`) {
		t.Fatalf("expected summary in lite json, got %s", out)
	}
	if strings.Contains(out, `"provider"`) {
		t.Fatalf("did not expect provider metadata in lite json, got %s", out)
	}
}

func TestPrintJSONRawKeepsRawFields(t *testing.T) {
	results := []trending.Result{
		{
			Provider:    "60s API",
			Platform:    trending.PlatformXiaohongshu,
			DisplayName: "Xiaohongshu",
			Endpoint:    "/v2/rednote",
			FetchedAt:   time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
			Items: []trending.Item{
				{
					Rank:  1,
					Title: "topic",
					Raw:   map[string]any{"score": "123w"},
				},
			},
		},
	}

	var buf bytes.Buffer
	PrintJSONRaw(&buf, results)
	out := buf.String()

	if !strings.Contains(out, `"provider": "60s API"`) {
		t.Fatalf("expected provider metadata in raw json, got %s", out)
	}
	if !strings.Contains(out, `"raw"`) {
		t.Fatalf("expected raw payload in raw json, got %s", out)
	}
}

func TestPrintDetailJSONIncludesSummaryButNotProviderMetadata(t *testing.T) {
	result := trending.Result{
		Provider:    "60s API",
		Platform:    trending.PlatformZhihu,
		DisplayName: "Zhihu",
		FetchedAt:   time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
		TotalCount:  30,
	}
	item := trending.Item{
		Rank:    1,
		Title:   "topic",
		Summary: "summary text",
		URL:     "https://example.com",
		Raw:     map[string]any{"detail": "raw detail"},
	}

	var buf bytes.Buffer
	PrintDetailJSON(&buf, result, item)
	out := buf.String()

	if !strings.Contains(out, `"summary": "summary text"`) {
		t.Fatalf("expected summary in detail json, got %s", out)
	}
	if strings.Contains(out, `"provider"`) {
		t.Fatalf("did not expect provider metadata in trimmed detail json, got %s", out)
	}
	if strings.Contains(out, `"raw"`) {
		t.Fatalf("did not expect raw payload in trimmed detail json, got %s", out)
	}
}

func TestPrintDetailJSONRawKeepsProviderAndRaw(t *testing.T) {
	result := trending.Result{
		Provider:    "60s API",
		Platform:    trending.PlatformZhihu,
		DisplayName: "Zhihu",
		Endpoint:    "/v2/zhihu",
		FetchedAt:   time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
		TotalCount:  30,
	}
	item := trending.Item{
		Rank:  1,
		Title: "topic",
		Raw:   map[string]any{"detail": "raw detail"},
	}

	var buf bytes.Buffer
	PrintDetailJSONRaw(&buf, result, item)
	out := buf.String()

	if !strings.Contains(out, `"provider": "60s API"`) {
		t.Fatalf("expected provider metadata in raw detail json, got %s", out)
	}
	if !strings.Contains(out, `"raw"`) {
		t.Fatalf("expected raw payload in raw detail json, got %s", out)
	}
}
