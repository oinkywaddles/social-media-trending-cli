---
name: social-media-trending
description: Query trending topics from major Chinese social/content platforms
tools: [Bash]
---

# Social Media Trending CLI

Use this CLI to query trending lists across:

- `xiaohongshu`
- `douyin`
- `bilibili`
- `weibo`
- `zhihu`
- `dongchedi`

Read:

- Setup/install issues -> `SETUP.md`
- Runtime behavior / output / snapshots -> `MECHANISM.md`
- Errors during use -> `TROUBLESHOOTING.md`

## Commands

### platforms

Show supported platforms and aliases.

```bash
social-media-trending-cli platforms
social-media-trending-cli platforms --json
```

### list

Fetch one or more trending lists.

```bash
# One platform
social-media-trending-cli list xiaohongshu
social-media-trending-cli list douyin --limit 10

# Multiple platforms
social-media-trending-cli list xiaohongshu douyin weibo
social-media-trending-cli list all

# Output formats
social-media-trending-cli list zhihu --json
social-media-trending-cli list zhihu --json-raw

# Freshness / request control
social-media-trending-cli list all --no-cache
social-media-trending-cli list all --cache-ttl 5m
social-media-trending-cli list all --min-interval 2s
social-media-trending-cli list all --timeout 20s
```

### detail

Show one ranked item in more detail.

```bash
social-media-trending-cli detail zhihu 1
social-media-trending-cli detail xiaohongshu 3 --json
social-media-trending-cli detail douyin 2 --json-raw
```

Use this after `list` when the user wants the item's:

- URL
- summary
- timestamps
- cover
- raw payload

## Usage guidelines

- Use plain text output when the user wants to scan rankings quickly
- Use `detail` instead of `list --json-raw` when only one ranked item matters
- Use `--json` for downstream summarization, clustering, or LLM consumption
- Use `--json-raw` only when you need provider metadata or original upstream fields
- Run `platforms` first if the user is unsure about supported platform names or aliases
