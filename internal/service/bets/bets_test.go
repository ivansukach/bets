package bets

import (
	"context"
	betsRps "github.com/ivansukach/bets/internal/repositories/bets"
	blockedBetsRps "github.com/ivansukach/bets/internal/repositories/blocked-bets"
	usersRps "github.com/ivansukach/bets/internal/repositories/blocked-users"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"testing"
	"time"
)

func TestCreateBet(t *testing.T) {
	db, err := sqlx.Connect("postgres",
		"user=su password=su "+
			"host=localhost dbname=auth")
	if err != nil {
		log.Fatal(err)
	}
	rpsOfBlockedUsers := usersRps.New(db)
	rpsOfBlockedBets := blockedBetsRps.New(db)
	rpsOfBets := betsRps.New(db)
	betsService := New(rpsOfBlockedUsers, rpsOfBlockedBets, rpsOfBets)
	id := int64(rand.Intn(int(time.Now().Unix())))
	err = betsService.CreateUserBetSummary(context.Background(), &betsRps.Bet{
		Id:                 id,
		Name:               "Isaak Newton",
		SumAmount:          200.22,
		AverageCoefficient: 1.143,
		NumberOfBets:       36,
	})
	if err != nil {
		log.Error(err)
	}
}
func TestCreateAndDeleteBlockedBet(t *testing.T) {
	db, err := sqlx.Connect("postgres",
		"user=su password=su "+
			"host=localhost dbname=bets")
	if err != nil {
		log.Fatal(err)
	}
	rpsOfBlockedUsers := usersRps.New(db)
	rpsOfBlockedBets := blockedBetsRps.New(db)
	rpsOfBets := betsRps.New(db)
	betsService := New(rpsOfBlockedUsers, rpsOfBlockedBets, rpsOfBets)
	id := int64(rand.Intn(int(time.Now().Unix() / 2)))
	err = betsService.CreateBlockedUserBetSummary(context.Background(), &blockedBetsRps.BlockedBet{
		Id:                 id,
		Name:               "Nikola Tesla",
		SumAmount:          800.78,
		AverageCoefficient: 8.857,
		NumberOfBets:       64,
	})
	if err != nil {
		log.Error(err)
	}
	err = betsService.betsRps.Delete(context.Background(), id)
	if err != nil {
		log.Error(err)
	}
}
