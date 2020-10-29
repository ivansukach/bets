package blocked_users

import "context"

type BlockedUser struct {
	Id int64 `db:"id"`
}
type Repository interface {
	Create(ctx context.Context, user *BlockedUser) error
	Get(ctx context.Context, id int64) (*BlockedUser, error)
	Delete(ctx context.Context, id int64) error
	Listing(ctx context.Context) ([]BlockedUser, error)
}
