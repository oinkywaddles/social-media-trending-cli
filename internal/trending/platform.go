package trending

import (
	"fmt"
	"strings"
)

type Platform string

const (
	PlatformXiaohongshu Platform = "xiaohongshu"
	PlatformDouyin      Platform = "douyin"
	PlatformBilibili    Platform = "bilibili"
	PlatformWeibo       Platform = "weibo"
	PlatformZhihu       Platform = "zhihu"
	PlatformDongchedi   Platform = "dongchedi"
)

type PlatformInfo struct {
	Name        Platform `json:"name"`
	DisplayName string   `json:"display_name"`
	Endpoint    string   `json:"endpoint"`
	Aliases     []string `json:"aliases"`
}

var supportedPlatforms = []PlatformInfo{
	{Name: PlatformXiaohongshu, DisplayName: "Xiaohongshu", Endpoint: "/v2/rednote", Aliases: []string{"xiaohongshu", "rednote", "xhs"}},
	{Name: PlatformDouyin, DisplayName: "Douyin", Endpoint: "/v2/douyin", Aliases: []string{"douyin", "dy"}},
	{Name: PlatformBilibili, DisplayName: "Bilibili", Endpoint: "/v2/bili", Aliases: []string{"bilibili", "bili", "b站"}},
	{Name: PlatformWeibo, DisplayName: "Weibo", Endpoint: "/v2/weibo", Aliases: []string{"weibo", "wb"}},
	{Name: PlatformZhihu, DisplayName: "Zhihu", Endpoint: "/v2/zhihu", Aliases: []string{"zhihu", "知乎"}},
	{Name: PlatformDongchedi, DisplayName: "Dongchedi", Endpoint: "/v2/dongchedi", Aliases: []string{"dongchedi", "dcd", "懂车帝"}},
}

func SupportedPlatformInfos() []PlatformInfo {
	out := make([]PlatformInfo, len(supportedPlatforms))
	copy(out, supportedPlatforms)
	return out
}

func ResolvePlatforms(args []string) ([]Platform, error) {
	if len(args) == 0 || (len(args) == 1 && strings.EqualFold(args[0], "all")) {
		return AllPlatforms(), nil
	}

	seen := make(map[Platform]bool)
	result := make([]Platform, 0, len(args))
	for _, arg := range args {
		platform, err := ParsePlatform(arg)
		if err != nil {
			return nil, err
		}
		if !seen[platform] {
			seen[platform] = true
			result = append(result, platform)
		}
	}
	return result, nil
}

func ParsePlatform(value string) (Platform, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	for _, info := range supportedPlatforms {
		for _, alias := range info.Aliases {
			if normalized == strings.ToLower(alias) {
				return info.Name, nil
			}
		}
	}
	return "", fmt.Errorf("unsupported platform %q", value)
}

func AllPlatforms() []Platform {
	result := make([]Platform, 0, len(supportedPlatforms))
	for _, info := range supportedPlatforms {
		result = append(result, info.Name)
	}
	return result
}

func (p Platform) Info() PlatformInfo {
	for _, info := range supportedPlatforms {
		if info.Name == p {
			return info
		}
	}
	return PlatformInfo{Name: p, DisplayName: string(p)}
}

func (p Platform) String() string {
	return string(p)
}
