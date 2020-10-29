package bets

import "context"

type BlockedBet struct {
	Id                 int64   `db:"id"`
	Name               string  `db:"name"`
	SumAmount          float64 `db:"sum_amount"`
	AverageCoefficient float64 `db:"average_coefficient"`
	NumberOfBets       int64   `db:"number_of_bets"`
}
type Repository interface {
	Create(ctx context.Context, bet *BlockedBet) error
	Get(ctx context.Context, id int64) (*BlockedBet, error)
	Update(ctx context.Context, bet *BlockedBet) error
	Delete(ctx context.Context, id int64) error
	Listing(ctx context.Context) ([]BlockedBet, error)
}
