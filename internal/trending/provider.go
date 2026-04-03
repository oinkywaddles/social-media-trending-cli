package trending

import "context"

type Provider interface {
	Name() string
	Fetch(ctx context.Context, platform Platform) (Result, error)
}
