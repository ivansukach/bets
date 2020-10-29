package bets

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type blockedBetsRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &blockedBetsRepository{db: db}
}

func (r *blockedBetsRepository) Create(ctx context.Context, bet *BlockedBet) error {
	_, err := r.db.NamedExec("INSERT INTO blocked_bets VALUES (:id, :name, :sum_amount, :average_coefficient, :number_of_bets)", bet)
	return err
}
func (r *blockedBetsRepository) Get(ctx context.Context, id int64) (*BlockedBet, error) {
	b := BlockedBet{}
	err := r.db.QueryRowx("SELECT * FROM blocked_bets WHERE Id=$1", id).StructScan(&b)
	return &b, err
}
func (r *blockedBetsRepository) Update(ctx context.Context, bet *BlockedBet) error {
	_, err := r.db.NamedExec("UPDATE blocked_bets SET Sum_amount=:sum_amount, "+
		"Average_coefficient=:average_coefficient, Number_of_bets=:number_of_bets WHERE Id=:id", bet)
	return err
}
func (r *blockedBetsRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec("DELETE FROM blocked_bets WHERE Id=$1", id)
	return err
}
func (r *blockedBetsRepository) Listing(ctx context.Context) ([]BlockedBet, error) {
	var amount int64
	err := r.db.Get(&amount, "SELECT count(*) FROM blocked_bets")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx("SELECT * FROM blocked_bets")
	if err != nil {
		return nil, err
	}
	b := make([]BlockedBet, 0, amount)
	for rows.Next() {
		bet := BlockedBet{}
		err = rows.StructScan(&bet)
		if err != nil {
			return nil, err
		}
		b = append(b, bet)
	}
	return b, err
}
