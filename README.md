# social-media-trending-cli

[![Go](https://img.shields.io/badge/go-1.26%2B-00ADD8.svg)](#installation)
[![Platforms](https://img.shields.io/badge/platforms-6-2ea44f.svg)](#supported-platforms)
[![Agent%20Skill](https://img.shields.io/badge/agent%20skill-included-ff9f1c.svg)](#use-as-an-ai-agent-skill)

A Go CLI for fetching trending topics from major Chinese social/content platforms, designed for both terminal users and AI agents.

[English](#english) | [中文](#chinese)

> **AI Agent Tip:** Prefer `--json` for structured downstream use. Use `detail <platform> <rank>` after `list` when you want one topic's URL, summary, or metadata without parsing the full raw payload.

## English

### What This Is

`social-media-trending-cli` gives you one consistent interface for hot/trending topics across:

- Xiaohongshu / Rednote
- Douyin
- Bilibili
- Weibo
- Zhihu
- Dongchedi

It is built for two workflows:

- **Human workflow**: scan current hot topics in the terminal
- **Agent workflow**: give an AI agent a stable command surface for trending data, instead of ad hoc scraping logic per platform

Current upstream provider:

- `60s API` (`https://60s.viki.moe`)

### Use As an AI Agent Skill

This repository ships with a local skill at:

- `skills/social-media-trending-cli/`

That skill exists to help agents fetch hot lists and item-level detail in a predictable way.

Included files:

- [`skills/social-media-trending-cli/SKILL.md`](./skills/social-media-trending-cli/SKILL.md) — how to use the CLI
- [`skills/social-media-trending-cli/SETUP.md`](./skills/social-media-trending-cli/SETUP.md) — install and verification
- [`skills/social-media-trending-cli/MECHANISM.md`](./skills/social-media-trending-cli/MECHANISM.md) — runtime behavior, snapshots, output modes
- [`skills/social-media-trending-cli/TROUBLESHOOTING.md`](./skills/social-media-trending-cli/TROUBLESHOOTING.md) — common failures and fixes

Use this skill when an agent needs to:

- fetch current hot topics from one or more platforms
- inspect one ranked item after a list command
- retrieve structured JSON for summarization, analysis, monitoring, or routing

### Features

- **Multi-platform trending** — one CLI for 6 platforms
- **Agent-friendly output** — default text, trimmed `--json`, full `--json-raw`
- **Detail lookup** — inspect one ranked item with `detail <platform> <rank>`
- **Snapshot-backed workflow** — `detail` can reuse the latest platform ranking
- **Conservative request handling** — local cache, retry, throttling
- **Simple terminal UX** — good defaults, readable tables, stable command names

### Supported Platforms

| Platform | Command Name | Common Aliases | Upstream Endpoint |
|----------|--------------|----------------|-------------------|
| Xiaohongshu | `xiaohongshu` | `rednote`, `xhs` | `/v2/rednote` |
| Douyin | `douyin` | `dy` | `/v2/douyin` |
| Bilibili | `bilibili` | `bili`, `b站` | `/v2/bili` |
| Weibo | `weibo` | `wb` | `/v2/weibo` |
| Zhihu | `zhihu` | `知乎` | `/v2/zhihu` |
| Dongchedi | `dongchedi` | `dcd`, `懂车帝` | `/v2/dongchedi` |

### Installation

Build from source:

```bash
git clone <your-repo-url>
cd social_media_trending_cli
go build -o social-media-trending-cli .
```

Or install into your Go bin:

```bash
git clone <your-repo-url>
cd social_media_trending_cli
go install .
```

Run directly from source if preferred:

```bash
go run . --help
```

### Quick Start

```bash
# List supported platforms
go run . platforms

# Fetch current rankings
go run . list all --limit 3
go run . list zhihu --limit 5

# Inspect one ranked item
go run . detail zhihu 1

# Structured output for programs / agents
go run . list zhihu --json
go run . detail douyin 2 --json
```

### Usage

#### `platforms`

```bash
social-media-trending-cli platforms
social-media-trending-cli platforms --json
```

#### `list`

```bash
# One platform
social-media-trending-cli list xiaohongshu
social-media-trending-cli list douyin --limit 10

# Multiple platforms
social-media-trending-cli list xiaohongshu douyin weibo
social-media-trending-cli list all

# Output modes
social-media-trending-cli list zhihu --json
social-media-trending-cli list zhihu --json-raw

# Freshness / request control
social-media-trending-cli list all --no-cache
social-media-trending-cli list all --cache-ttl 5m
social-media-trending-cli list all --min-interval 2s
social-media-trending-cli list all --timeout 20s
```

#### `detail`

```bash
social-media-trending-cli detail zhihu 1
social-media-trending-cli detail xiaohongshu 3 --json
social-media-trending-cli detail douyin 2 --json-raw
```

Use `detail` when a topic from `list` looks interesting and you want:

- the URL
- summary / description
- timestamps
- cover
- raw upstream fields

### Output Modes

Default output:

- `list` -> Markdown table
- `detail` -> structured text

`--json`:

- trimmed for downstream programs and LLMs
- keeps user-relevant fields such as `title`, `score`, `url`, `summary`

`--json-raw`:

- full normalized structure
- includes provider metadata and original upstream `raw`

### Snapshot Behavior

`detail <platform> <rank>` uses the latest platform snapshot when available.

- snapshots refresh whenever `list` runs
- snapshot TTL is `1h`
- if no fresh snapshot exists, `detail` fetches live data and saves a new snapshot
- `--limit` only affects displayed output; the full ranking is still saved for later `detail` lookup

For deeper runtime notes, see [`MECHANISM.md`](./skills/social-media-trending-cli/MECHANISM.md).

### Project Layout

```text
social_media_trending_cli/
├── cmd/                                  # Cobra commands
├── internal/
│   ├── cache/                            # HTTP cache
│   ├── httpx/                            # HTTP client, retry, throttle
│   ├── output/                           # Text / JSON renderers
│   ├── snapshot/                         # Latest per-platform snapshot store
│   └── trending/                         # Platform model + provider adapters
├── skills/
│   └── social-media-trending-cli/        # Agent skill docs
│       ├── SKILL.md
│       ├── SETUP.md
│       ├── MECHANISM.md
│       └── TROUBLESHOOTING.md
├── README.md
└── main.go
```

## Chinese

### 这是什么

`social-media-trending-cli` 是一个用于获取中文社交/内容平台热榜的 Go CLI，同时面向两类使用场景：

- 给人直接在终端里看热榜
- 给 Agent 提供稳定的热榜数据入口

当前支持：

- 小红书
- 抖音
- B 站
- 微博
- 知乎
- 懂车帝

### 这个仓库里的 skill 在哪里

本仓库自带一个本地 skill，路径是：

- `skills/social-media-trending-cli/`

它的目的不是讲实现细节，而是教 Agent 如何使用这个 CLI 来获取热榜数据、查看某一条热点详情，以及输出结构化 JSON。

主要文件：

- `SKILL.md`：怎么用
- `SETUP.md`：安装与验证
- `MECHANISM.md`：机制说明
- `TROUBLESHOOTING.md`：问题排查

### 常用命令

```bash
# 查看支持的平台
go run . platforms

# 拉热榜
go run . list all --limit 3
go run . list zhihu --json

# 看某一条热点详情
go run . detail zhihu 1
go run . detail xiaohongshu 3 --json
```

### 适合 Agent 的原因

- 命令稳定，平台统一
- 默认文本输出适合直接看
- `--json` 适合摘要、分析、路由
- `detail` 可以从最新快照里取某一条热点，避免重复手动解析榜单

### 更多说明

如果你关心输出结构、快照机制、缓存、限频和上游行为，请看：

- [`skills/social-media-trending-cli/MECHANISM.md`](./skills/social-media-trending-cli/MECHANISM.md)
