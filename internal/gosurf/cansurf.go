package gosurf

import (
	"context"
)

func CanSurf(ctx context.Context, spot Spot) (bool, error) {
	return true, nil
}
