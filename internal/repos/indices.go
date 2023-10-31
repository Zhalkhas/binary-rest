package repos

import "context"

type Indices interface {
	Search(ctx context.Context, val int) (int, error)
}
