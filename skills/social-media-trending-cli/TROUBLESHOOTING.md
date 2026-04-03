# Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `unsupported platform "..."` | Platform name or alias is invalid | Run `social-media-trending-cli platforms` and use one of the supported names |
| `rank must be a positive integer` | `detail` rank arg is missing or invalid | Use `detail <platform> <rank>` where rank is `1`, `2`, `3`, ... |
| `... rank N not found; available range is 1-M` | Requested rank exceeds the current snapshot/live ranking size | Use a smaller rank or run `list <platform>` first to inspect available items |
| `use either --json or --json-raw, not both` | Mutually exclusive output flags were passed together | Choose exactly one of `--json` or `--json-raw` |
| `GET https://60s.viki.moe/... returned 429` | Upstream is rate limiting or temporarily overloaded | Wait and retry; increase `--min-interval`; avoid repeated `list all --no-cache` runs |
| `GET https://60s.viki.moe/... returned 5xx` | Upstream service failure | Retry later; the client already retries transient failures up to 3 times |
| `context deadline exceeded` | Network is slow or upstream did not respond before timeout | Increase `--timeout`, check network connectivity, retry |
| `command not found: social-media-trending-cli` | Binary was not installed or not in `PATH` | Use `go run . ...` from the repo, or build/install the binary first |
| `failed to read snapshot` / `failed to save snapshot` | Snapshot directory is unavailable or not writable | Check permissions under the user cache directory; rerun with a writable home/cache dir |
| `failed to read cache` / stale rankings appear unexpected | HTTP cache is still fresh | Use `--no-cache` or lower `--cache-ttl` when you need a fresh upstream fetch |

## Notes on rate limits

The current upstream (`60s API`) does not publish a clear quota policy in the CLI's verified sources.

Because of that, the CLI already applies conservative protection:

- sequential fetches
- request throttling
- local HTTP cache
- short retry with backoff

If you still see upstream instability, reduce request frequency rather than increasing parallelism.
