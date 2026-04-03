package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"social-media-trending-cli/internal/cache"
)

const defaultUserAgent = "social-media-trending-cli/0.1"
const maxAttempts = 3

type Client struct {
	base        *http.Client
	cache       *cache.FileCache
	minInterval time.Duration
	userAgent   string

	mu            sync.Mutex
	nextAvailable time.Time
}

func NewClient(timeout, minInterval time.Duration, fileCache *cache.FileCache) *Client {
	return &Client{
		base: &http.Client{
			Timeout: timeout,
		},
		cache:       fileCache,
		minInterval: minInterval,
		userAgent:   defaultUserAgent,
	}
}

func (c *Client) GetJSON(ctx context.Context, url, cacheKey string, ttl time.Duration, target any) (bool, time.Time, error) {
	if c.cache != nil && cacheKey != "" && ttl > 0 {
		body, fetchedAt, fresh, err := c.cache.ReadFresh(cacheKey, ttl)
		if err == nil && fresh {
			if err := json.Unmarshal(body, target); err != nil {
				return false, time.Time{}, err
			}
			return true, fetchedAt, nil
		}
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.wait(ctx); err != nil {
			return false, time.Time{}, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return false, time.Time{}, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.base.Do(req)
		if err != nil {
			lastErr = err
			if attempt == maxAttempts {
				break
			}
			if err := sleepContext(ctx, retryDelay(attempt)); err != nil {
				return false, time.Time{}, err
			}
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			if attempt == maxAttempts {
				break
			}
			if err := sleepContext(ctx, retryDelay(attempt)); err != nil {
				return false, time.Time{}, err
			}
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("GET %s returned %d: %s", url, resp.StatusCode, shortenBody(string(body)))
			if attempt == maxAttempts || !retryableStatus(resp.StatusCode) {
				break
			}
			if err := sleepContext(ctx, retryDelay(attempt)); err != nil {
				return false, time.Time{}, err
			}
			continue
		}

		fetchedAt := time.Now().UTC()
		if c.cache != nil && cacheKey != "" && ttl > 0 {
			_ = c.cache.Write(cacheKey, body, fetchedAt)
		}

		if err := json.Unmarshal(body, target); err != nil {
			return false, time.Time{}, err
		}
		return false, fetchedAt, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("GET %s failed after %d attempts", url, maxAttempts)
	}
	return false, time.Time{}, lastErr
}

func (c *Client) wait(ctx context.Context) error {
	if c.minInterval <= 0 {
		return nil
	}

	c.mu.Lock()
	now := time.Now()
	waitFor := time.Duration(0)
	if now.Before(c.nextAvailable) {
		waitFor = c.nextAvailable.Sub(now)
		now = c.nextAvailable
	}
	c.nextAvailable = now.Add(c.minInterval)
	c.mu.Unlock()

	if waitFor <= 0 {
		return nil
	}

	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func shortenBody(body string) string {
	body = strings.TrimSpace(body)
	if len(body) <= 160 {
		return body
	}
	return body[:160] + "..."
}

func retryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= 500
}

func retryDelay(attempt int) time.Duration {
	return time.Duration(attempt*attempt) * 300 * time.Millisecond
}

func sleepContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
