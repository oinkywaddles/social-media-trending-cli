# Setup Guide

## Install

Build from source:

```bash
git clone <your-repo-url>
cd social_media_trending_cli
go build -o social-media-trending-cli .
```

Or install into your Go bin directory:

```bash
git clone <your-repo-url>
cd social_media_trending_cli
go install .
```

Then verify:

```bash
social-media-trending-cli --help
```

If you are running from source directly, use:

```bash
go run . --help
```

## Configuration

No API key is required.

For runtime behavior, cache/snapshot rules, output formats, and upstream details, read `MECHANISM.md`.

## Verify

```bash
social-media-trending-cli platforms
social-media-trending-cli list zhihu --limit 2
social-media-trending-cli detail zhihu 1 --json
```

If you encounter errors, read `skills/TROUBLESHOOTING.md`.
