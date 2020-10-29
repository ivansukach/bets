package blocking

import (
	"context"
	usersRps "github.com/ivansukach/bets/internal/repositories/blocked-users"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"testing"
	"time"
)

func TestBlockingService(t *testing.T) {
	db, err := sqlx.Connect("postgres",
		"user=su password=su "+
			"host=localhost dbname=auth")
	if err != nil {
		log.Fatal(err)
	}
	rpsOfBlockedUsers := usersRps.New(db)
	blockingService := New(rpsOfBlockedUsers)
	_, err = blockingService.BlockedReport(context.Background())
	if err == nil {
		log.Error("Another db. Table does not exist. Error should be not nil ")
	}
}
func TestCreateAndDeleteUser(t *testing.T) {
	db, err := sqlx.Connect("postgres",
		"user=su password=su "+
			"host=localhost dbname=bets")
	if err != nil {
		log.Fatal(err)
	}
	rpsOfBlockedUsers := usersRps.New(db)
	blockingService := New(rpsOfBlockedUsers)
	id := int64(rand.Intn(int(time.Now().Unix())))
	err = blockingService.BlockUsers(context.Background(), []int64{id})
	if err != nil {
		log.Error(err)
	}
	err = blockingService.DeleteUser(context.Background(), id)
	if err != nil {
		log.Error(err)
	}
}
