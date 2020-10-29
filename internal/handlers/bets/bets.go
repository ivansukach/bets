package bets

import (
	"github.com/ivansukach/bets/internal/service/bets"
)

type Bets struct {
	betsService *bets.Service
}

func New(betsService *bets.Service) *Bets {
	return &Bets{betsService: betsService}
}
