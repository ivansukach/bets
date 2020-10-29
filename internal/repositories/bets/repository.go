package bets

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type betsRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &betsRepository{db: db}
}

func (r *betsRepository) Create(ctx context.Context, bet *Bet) error {
	_, err := r.db.NamedExec("INSERT INTO bets VALUES (:id, :name, :sum_amount, :average_coefficient, :number_of_bets)", bet)
	return err
}
func (r *betsRepository) Get(ctx context.Context, id int64) (*Bet, error) {
	b := Bet{}
	err := r.db.QueryRowx("SELECT * FROM bets WHERE Id=$1", id).StructScan(&b)
	return &b, err
}
func (r *betsRepository) Update(ctx context.Context, bet *Bet) error {
	_, err := r.db.NamedExec("UPDATE bets SET Sum_amount=:sum_amount, "+
		"Average_coefficient=:average_coefficient, Number_of_bets=:number_of_bets WHERE Id=:id", bet)
	return err
}
func (r *betsRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec("DELETE FROM bets WHERE Id=$1", id)
	return err
}
func (r *betsRepository) Listing(ctx context.Context) ([]Bet, error) {
	var amount int64
	err := r.db.Get(&amount, "SELECT count(*) FROM bets")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx("SELECT * FROM bets")
	if err != nil {
		return nil, err
	}
	b := make([]Bet, 0, amount)
	for rows.Next() {
		bet := Bet{}
		err = rows.StructScan(&bet)
		if err != nil {
			return nil, err
		}
		b = append(b, bet)
	}
	return b, err
}
