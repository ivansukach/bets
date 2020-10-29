package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ivansukach/bets/internal/config"
	"github.com/ivansukach/bets/internal/handlers/bets"
	"github.com/ivansukach/bets/internal/handlers/blocking"
	betsRps "github.com/ivansukach/bets/internal/repositories/bets"
	blockedBetsRps "github.com/ivansukach/bets/internal/repositories/blocked-bets"
	usersRps "github.com/ivansukach/bets/internal/repositories/blocked-users"
	betsSrv "github.com/ivansukach/bets/internal/service/bets"
	blockingSrv "github.com/ivansukach/bets/internal/service/blocking"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	log.SetLevel(log.DebugLevel)
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}
func GetRouter() http.Handler {
	db, err := sqlx.Connect("postgres",
		"user=su password=su "+
			"host=localhost dbname=bets")
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	rpsOfBlockedUsers := usersRps.New(db)
	rpsOfBlockedBets := blockedBetsRps.New(db)
	rpsOfBets := betsRps.New(db)
	blockingService := blockingSrv.New(rpsOfBlockedUsers)
	betsService := betsSrv.New(rpsOfBlockedUsers, rpsOfBlockedBets, rpsOfBets)
	blockingHandlers := blocking.New(blockingService)
	betsHandlers := bets.New(betsService)
	r.Use(middleware.Logger)
	r.Post("/block-users", blockingHandlers.BlockUsers)
	r.Post("/process", betsHandlers.Process)
	r.Get("/report", betsHandlers.Report)
	r.Get("/blocked-report", blockingHandlers.BlockedReport)
	return r
}

func main() {
	log.Info("Client started")
	cfg := config.Load()
	log.Fatal(http.ListenAndServe(":"+cfg.Port, GetRouter()))
}
