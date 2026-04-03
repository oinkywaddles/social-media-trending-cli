# Mechanism Guide

## Upstream source

Current upstream provider:

- `60s API` at `https://60s.viki.moe`

Supported endpoints currently used by this CLI:

- `xiaohongshu` -> `/v2/rednote`
- `douyin` -> `/v2/douyin`
- `bilibili` -> `/v2/bili`
- `weibo` -> `/v2/weibo`
- `zhihu` -> `/v2/zhihu`
- `dongchedi` -> `/v2/dongchedi`

## Output modes

Default output:

- `list` -> human-readable Markdown table
- `detail` -> structured text

`--json`:

- trimmed output for downstream programs or LLMs
- keeps user-relevant fields like `title`, `score`, `url`, `summary`
- omits provider metadata and upstream `raw`

`--json-raw`:

- full normalized structure
- includes provider metadata, request policy, and original upstream `raw`

Error JSON format:

```json
{"error":"error message here"}
```

## Snapshot behavior

`detail <platform> <rank>` uses a per-platform latest snapshot.

- snapshots are refreshed whenever `list` runs
- snapshot TTL is `1h`
- if no fresh snapshot exists, `detail` fetches the platform live, saves a new snapshot, then returns the ranked item
- snapshots are stored per platform, not per invocation
- `--limit` only changes displayed output; the full platform ranking is still saved for later `detail` lookup

Typical macOS snapshot paths:

- `~/Library/Caches/social-media-trending-cli/snapshots/zhihu.json`
- `~/Library/Caches/social-media-trending-cli/snapshots/xiaohongshu.json`

## HTTP cache and request behavior

No API key is required.

Client-side protections are enabled by default:

- HTTP timeout: `15s`
- minimum interval between upstream requests: `1.2s`
- HTTP response cache TTL: `2m`
- transient failure retry count: `3`
- multi-platform fetches are sequential, not parallel

Typical macOS HTTP cache path:

- `~/Library/Caches/social-media-trending-cli/`

## API cost

Approximate upstream call count:

- `platforms`: `0`
- `list <platform>`: `1`
- `list all`: up to `6`
- `detail <platform> <rank>`: `0` if a fresh snapshot exists, otherwise `1`

## JSON shapes

### list --json

```json
{
  "generated_at": "2026-04-04T00:00:00Z",
  "results": [
    {
      "platform": "zhihu",
      "fetched_at": "2026-04-04T00:00:00Z",
      "total_count": 30,
      "returned_count": 2,
      "items": [
        {
          "rank": 1,
          "title": "topic title",
          "score": "539 万热度",
          "meta": "239 answers | 564 followers",
          "url": "https://www.zhihu.com/question/...",
          "summary": "summary text",
          "published_at": "2026/04/03 13:29:58"
        }
      ]
    }
  ]
}
```

### detail --json

```json
{
  "generated_at": "2026-04-04T00:00:00Z",
  "platform": "douyin",
  "fetched_at": "2026-04-04T00:00:00Z",
  "total_count": 49,
  "item": {
    "rank": 2,
    "title": "topic title",
    "score": "11,909,687",
    "meta": "active 2026-04-04 00:37:05",
    "url": "https://www.douyin.com/search/...",
    "cover": "https://...",
    "published_at": "2026/04/03 15:14:21",
    "updated_at": "2026-04-04 00:37:05"
  }
}
```

## Notes on rate limits

The current upstream does not publish a clear quota policy in the CLI's verified sources.

Because of that, the CLI already applies conservative protection:

- request throttling
- local HTTP cache
- retry with backoff
- sequential multi-platform fetches
