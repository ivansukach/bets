package blocked_users

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type blockedUsersRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &blockedUsersRepository{db: db}
}

func (r *blockedUsersRepository) Create(ctx context.Context, user *BlockedUser) error {
	_, err := r.db.NamedExec("INSERT INTO blocked_users VALUES (:id)", user)
	return err
}
func (r *blockedUsersRepository) Get(ctx context.Context, id int64) (*BlockedUser, error) {
	u := BlockedUser{}
	err := r.db.QueryRowx("SELECT * FROM blocked_users WHERE Id=$1", id).StructScan(&u)
	return &u, err
}
func (r *blockedUsersRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec("DELETE FROM blocked_users WHERE id=$1", id)
	return err
}
func (r *blockedUsersRepository) Listing(ctx context.Context) ([]BlockedUser, error) {
	var amount int64
	err := r.db.Get(&amount, "SELECT count(*) FROM blocked_users")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx("SELECT * FROM blocked_users")
	if err != nil {
		return nil, err
	}
	u := make([]BlockedUser, 0, amount)
	for rows.Next() {
		usr := BlockedUser{}
		err = rows.StructScan(&usr)
		if err != nil {
			return nil, err
		}
		u = append(u, usr)
	}
	return u, err
}
