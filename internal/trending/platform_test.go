package trending

import "testing"

func TestResolvePlatformsDefaultsToAll(t *testing.T) {
	got, err := ResolvePlatforms(nil)
	if err != nil {
		t.Fatalf("ResolvePlatforms returned error: %v", err)
	}

	want := AllPlatforms()
	if len(got) != len(want) {
		t.Fatalf("expected %d platforms, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("platform %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestParsePlatformAliases(t *testing.T) {
	tests := map[string]Platform{
		"xhs":     PlatformXiaohongshu,
		"rednote": PlatformXiaohongshu,
		"dy":      PlatformDouyin,
		"bili":    PlatformBilibili,
		"b站":      PlatformBilibili,
		"wb":      PlatformWeibo,
		"知乎":      PlatformZhihu,
		"懂车帝":     PlatformDongchedi,
	}

	for input, want := range tests {
		got, err := ParsePlatform(input)
		if err != nil {
			t.Fatalf("ParsePlatform(%q) returned error: %v", input, err)
		}
		if got != want {
			t.Fatalf("ParsePlatform(%q): expected %q, got %q", input, want, got)
		}
	}
}
