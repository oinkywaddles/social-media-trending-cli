package trending

import "testing"

func TestNormalizeSixtySItemXiaohongshu(t *testing.T) {
	raw := map[string]any{
		"rank":           float64(2),
		"title":          "topic",
		"score":          "904.2w",
		"word_type":      "热",
		"work_type_icon": "https://img.example/icon.png",
		"link":           "https://www.xiaohongshu.com/search_result?keyword=topic",
	}

	item := normalizeSixtySItem(PlatformXiaohongshu, 0, raw)

	if item.Rank != 2 {
		t.Fatalf("expected rank 2, got %d", item.Rank)
	}
	if item.Score != "904.2w" {
		t.Fatalf("expected score 904.2w, got %q", item.Score)
	}
	if item.Label != "热" {
		t.Fatalf("expected label 热, got %q", item.Label)
	}
	if item.Cover != "https://img.example/icon.png" {
		t.Fatalf("expected cover to be carried over, got %q", item.Cover)
	}
}

func TestNormalizeSixtySItemDouyin(t *testing.T) {
	raw := map[string]any{
		"title":       "hot topic",
		"hot_value":   float64(12040689),
		"event_time":  "2026/04/02 10:56:49",
		"active_time": "2026-04-04 00:02:09",
		"link":        "https://www.douyin.com/search/topic",
	}

	item := normalizeSixtySItem(PlatformDouyin, 0, raw)

	if item.Rank != 1 {
		t.Fatalf("expected fallback rank 1, got %d", item.Rank)
	}
	if item.Score != "12,040,689" {
		t.Fatalf("expected formatted score, got %q", item.Score)
	}
	if item.Meta != "active 2026-04-04 00:02:09" {
		t.Fatalf("unexpected meta: %q", item.Meta)
	}
	if item.ScoreValue == nil || *item.ScoreValue != 12040689 {
		t.Fatalf("expected numeric score value 12040689, got %+v", item.ScoreValue)
	}
}

func TestNormalizeSixtySItemZhihu(t *testing.T) {
	raw := map[string]any{
		"title":          "question title",
		"detail":         "question detail",
		"hot_value_desc": "501 万热度",
		"answer_cnt":     float64(229),
		"follower_cnt":   float64(531),
		"comment_cnt":    float64(0),
		"created":        "2026/04/03 13:29:58",
		"link":           "https://www.zhihu.com/question/123",
	}

	item := normalizeSixtySItem(PlatformZhihu, 0, raw)

	if item.Score != "501 万热度" {
		t.Fatalf("expected zhihu score description, got %q", item.Score)
	}
	if item.Summary != "question detail" {
		t.Fatalf("expected summary to come from detail, got %q", item.Summary)
	}
	if item.Meta != "229 answers | 531 followers" {
		t.Fatalf("unexpected zhihu meta: %q", item.Meta)
	}
	if item.PublishedAt != "2026/04/03 13:29:58" {
		t.Fatalf("unexpected published time: %q", item.PublishedAt)
	}
}

func TestNormalizeSixtySItemDongchedi(t *testing.T) {
	raw := map[string]any{
		"rank":       float64(1),
		"title":      "car topic",
		"url":        "https://www.dongchedi.com/search?keyword=car",
		"score":      float64(963896),
		"score_desc": "96.4w",
	}

	item := normalizeSixtySItem(PlatformDongchedi, 0, raw)

	if item.Score != "96.4w" {
		t.Fatalf("expected score_desc to win, got %q", item.Score)
	}
	if item.ScoreValue == nil || *item.ScoreValue != 963896 {
		t.Fatalf("expected numeric score value 963896, got %+v", item.ScoreValue)
	}
}
