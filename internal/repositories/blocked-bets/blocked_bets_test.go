package bets

import (
	"context"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"testing"
	"time"
)

func TestCreateDeleteAndListingBlockedBet(t *testing.T) {
	db, err := sqlx.Connect("postgres",
		"user=su password=su "+
			"host=localhost dbname=bets")
	if err != nil {
		log.Fatal(err)
	}
	rpsOfBlockedUsers := New(db)
	id := int64(rand.Intn(int(time.Now().Unix())))
	err = rpsOfBlockedUsers.Create(context.Background(), &BlockedBet{
		Id:                 id,
		Name:               "Ivan Sukach",
		SumAmount:          137.3,
		AverageCoefficient: 2.53,
		NumberOfBets:       917,
	})
	if err != nil {
		log.Error(err)
	}
	err = rpsOfBlockedUsers.Delete(context.Background(), id)
	if err != nil {
		log.Error(err)
	}
	_, err = rpsOfBlockedUsers.Listing(context.Background())
	if err != nil {
		log.Error(err)
	}
}
