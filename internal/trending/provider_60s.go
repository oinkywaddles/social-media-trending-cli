package trending

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"social-media-trending-cli/internal/httpx"
)

const sixtySBaseURL = "https://60s.viki.moe"

type sixtySResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []map[string]any `json:"data"`
}

type SixtySProvider struct {
	client       *httpx.Client
	cacheTTL     time.Duration
	minInterval  time.Duration
	upstreamNote string
}

func NewSixtySProvider(client *httpx.Client, cacheTTL, minInterval time.Duration) *SixtySProvider {
	return &SixtySProvider{
		client:       client,
		cacheTTL:     cacheTTL,
		minInterval:  minInterval,
		upstreamNote: "No public rate-limit details found in 60s docs or response headers; client applies a conservative throttle",
	}
}

func (p *SixtySProvider) Name() string {
	return "60s API"
}

func (p *SixtySProvider) Fetch(ctx context.Context, platform Platform) (Result, error) {
	info := platform.Info()
	if info.Endpoint == "" {
		return Result{}, fmt.Errorf("unsupported platform %q", platform)
	}

	var resp sixtySResponse
	url := sixtySBaseURL + info.Endpoint
	cacheKey := "60s:" + platform.String()
	cached, fetchedAt, err := p.client.GetJSON(ctx, url, cacheKey, p.cacheTTL, &resp)
	if err != nil {
		return Result{}, err
	}
	if resp.Code != 200 {
		return Result{}, fmt.Errorf("upstream returned code %d: %s", resp.Code, resp.Message)
	}

	items := make([]Item, 0, len(resp.Data))
	for idx, raw := range resp.Data {
		items = append(items, normalizeSixtySItem(platform, idx, raw))
	}

	cacheTTLText := p.cacheTTL.String()
	if p.cacheTTL <= 0 {
		cacheTTLText = "disabled"
	}

	return Result{
		Provider:        p.Name(),
		Platform:        platform,
		DisplayName:     info.DisplayName,
		Endpoint:        info.Endpoint,
		FetchedAt:       fetchedAt,
		Cached:          cached,
		TotalCount:      len(items),
		UpstreamMessage: resp.Message,
		RequestPolicy: RequestPolicy{
			Upstream:         p.upstreamNote,
			LocalMinInterval: p.minInterval.String(),
			CacheTTL:         cacheTTLText,
		},
		Items: items,
	}, nil
}

func normalizeSixtySItem(platform Platform, index int, raw map[string]any) Item {
	item := Item{
		Rank: index + 1,
		Raw:  cloneMap(raw),
	}

	if rank, ok := int64Value(raw["rank"]); ok && rank > 0 {
		item.Rank = int(rank)
	}

	item.Title = stringValue(raw["title"])
	item.URL = firstNonEmpty(stringValue(raw["url"]), stringValue(raw["link"]))
	item.Cover = stringValue(raw["cover"])
	item.Label = firstNonEmpty(stringValue(raw["word_type"]), stringValue(raw["tag"]))
	item.Summary = firstNonEmpty(stringValue(raw["detail"]), stringValue(raw["desc"]))
	item.PublishedAt = firstNonEmpty(stringValue(raw["created"]), stringValue(raw["event_time"]))
	item.UpdatedAt = stringValue(raw["active_time"])

	switch platform {
	case PlatformXiaohongshu:
		item.Score = stringValue(raw["score"])
		item.Cover = firstNonEmpty(item.Cover, stringValue(raw["work_type_icon"]))
	case PlatformDouyin:
		setNumericScore(&item, raw["hot_value"])
		if item.UpdatedAt != "" {
			item.Meta = "active " + item.UpdatedAt
		} else if item.PublishedAt != "" {
			item.Meta = "event " + item.PublishedAt
		}
	case PlatformBilibili:
		// No score field in current upstream payload.
	case PlatformWeibo:
		setNumericScore(&item, raw["hot_value"])
	case PlatformZhihu:
		item.Score = stringValue(raw["hot_value_desc"])
		item.Meta = joinNonEmpty(
			countLabel(raw["answer_cnt"], "answers"),
			countLabel(raw["follower_cnt"], "followers"),
			countLabel(raw["comment_cnt"], "comments"),
		)
	case PlatformDongchedi:
		item.Score = firstNonEmpty(stringValue(raw["score_desc"]), stringValue(raw["score"]))
		if score, ok := int64Value(raw["score"]); ok {
			item.ScoreValue = &score
		}
	}

	return item
}

func setNumericScore(item *Item, value any) {
	if score, ok := int64Value(value); ok {
		item.ScoreValue = &score
		item.Score = formatInt(score)
		return
	}
	item.Score = stringValue(value)
}

func countLabel(value any, label string) string {
	count, ok := int64Value(value)
	if !ok || count <= 0 {
		return ""
	}
	return fmt.Sprintf("%s %s", formatInt(count), label)
}

func cloneMap(raw map[string]any) map[string]any {
	if raw == nil {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return raw
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return raw
	}
	return out
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case fmt.Stringer:
		return strings.TrimSpace(v.String())
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return ""
	}
}

func int64Value(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	case json.Number:
		n, err := v.Int64()
		return n, err == nil
	case string:
		cleaned := strings.TrimSpace(strings.ReplaceAll(v, ",", ""))
		if cleaned == "" {
			return 0, false
		}
		n, err := strconv.ParseInt(cleaned, 10, 64)
		return n, err == nil
	default:
		return 0, false
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
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

func formatInt(value int64) string {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	raw := strconv.FormatInt(value, 10)
	if len(raw) <= 3 {
		return sign + raw
	}

	var builder strings.Builder
	prefix := len(raw) % 3
	if prefix == 0 {
		prefix = 3
	}
	builder.WriteString(sign)
	builder.WriteString(raw[:prefix])
	for i := prefix; i < len(raw); i += 3 {
		builder.WriteByte(',')
		builder.WriteString(raw[i : i+3])
	}
	return builder.String()
}
