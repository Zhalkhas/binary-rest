package repos

import "context"

// Indices repository is responsible for searching indices of given int val
type Indices interface {
	Search(ctx context.Context, val int64) (int, error)
}
