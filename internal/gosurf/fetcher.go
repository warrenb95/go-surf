package gosurf

import "context"

type Fetcher interface {
	Get(ctx context.Context, spot Spot, params string) (map[string][]*Spot, error)
}
