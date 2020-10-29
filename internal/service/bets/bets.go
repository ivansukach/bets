package bets

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ivansukach/bets/internal/repositories/bets"
	blockedBets "github.com/ivansukach/bets/internal/repositories/blocked-bets"
	blockedUsers "github.com/ivansukach/bets/internal/repositories/blocked-users"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"strings"
	"sync"
)

type Service struct {
	blockedUsersRps blockedUsers.Repository
	blockedBetsRps  blockedBets.Repository
	betsRps         bets.Repository
	processing      bool
	mu              sync.Mutex
	//chanBets        chan map[int64]BetInfoSummary
	//chanBlockedBets chan map[int64]BetInfoSummary
}
type BetHistory struct {
	Name         string
	Bets         []float64
	Coefficients []float64
}
type BetInfo struct {
	Id          int64
	Name        string
	Amount      float64
	Coefficient float64
}
type BetInfoSummary struct {
	Id                 int64
	Name               string
	SumAmount          float64
	AverageCoefficient float64
}

func New(blockedUsersRps blockedUsers.Repository,
	blockedBetsRps blockedBets.Repository, betsRps bets.Repository) *Service {
	return &Service{blockedUsersRps: blockedUsersRps, blockedBetsRps: blockedBetsRps, betsRps: betsRps}
	//chanBets: make(chan map[int64]BetInfoSummary), chanBlockedBets: make(chan map[int64]BetInfoSummary

}
func (s *Service) FinishProcessing() {
	s.processing = false
}
func (s *Service) Process(ctx context.Context, file *multipart.File) error {
	s.processing = true
	s.mu.Lock()
	defer s.FinishProcessing()
	defer s.mu.Unlock()
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, *file)
	if err != nil {
		return err
	}
	records := strings.Split(string(buf.Bytes()), "\n")
	if records[len(records)-1] == "" {
		records = records[:len(records)-1]
	}
	log.Debug("Amount of records: ", len(records))
	betHistory := make(map[int64]BetHistory)
	blockedBetHistory := make(map[int64]BetHistory)
	for i := range records {
		betInfo := BetInfo{}
		values := strings.Split(records[i], ", ")
		fmt.Sscanf(values[0], "%d", &betInfo.Id)
		fmt.Sscanf(values[1], "%s", &betInfo.Name)
		fmt.Sscanf(values[2], "%f", &betInfo.Amount)
		fmt.Sscanf(values[3], "%f", &betInfo.Coefficient)
		//fmt.Sscanf(records[i], "%d, %s, %f, %f", &betInfo.Id, &betInfo.Name, &betInfo.Amount, &betInfo.Coefficient)
		if !s.BlockedUser(ctx, betInfo.Id) {
			_, exist := betHistory[betInfo.Id]
			if exist {
				betHistory[betInfo.Id] = BetHistory{
					Name:         betHistory[betInfo.Id].Name,
					Bets:         append(betHistory[betInfo.Id].Bets, betInfo.Amount),
					Coefficients: append(betHistory[betInfo.Id].Coefficients, betInfo.Coefficient),
				}
			} else {
				betHistory[betInfo.Id] = BetHistory{
					Name:         betInfo.Name,
					Bets:         []float64{betInfo.Amount},
					Coefficients: []float64{betInfo.Coefficient},
				}
			}
		} else {
			_, exist := blockedBetHistory[betInfo.Id]
			if exist {
				blockedBetHistory[betInfo.Id] = BetHistory{
					Name:         blockedBetHistory[betInfo.Id].Name,
					Bets:         append(blockedBetHistory[betInfo.Id].Bets, betInfo.Amount),
					Coefficients: append(blockedBetHistory[betInfo.Id].Coefficients, betInfo.Coefficient),
				}
			} else {
				blockedBetHistory[betInfo.Id] = BetHistory{
					Name:         betInfo.Name,
					Bets:         []float64{betInfo.Amount},
					Coefficients: []float64{betInfo.Coefficient},
				}
			}
		}
	}
	err = s.ComputeAverageValues(ctx, betHistory, blockedBetHistory)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) BlockedUser(ctx context.Context, id int64) bool {
	user, err := s.blockedUsersRps.Get(ctx, id)
	if user == nil || err != nil {
		return false
	}
	return true
}
func (s *Service) ComputeAverageValues(ctx context.Context, betHistory map[int64]BetHistory,
	blockedBetHistory map[int64]BetHistory) error {
	var err error
	for k, v := range betHistory {
		sumAmount := 0.0
		sumCoefficients := 0.0
		for _, value := range v.Bets {
			sumAmount += value
		}
		for _, value := range v.Coefficients {
			sumCoefficients += value
		}
		averageCoefficient := sumCoefficients / float64(len(v.Coefficients))
		bet := &bets.Bet{
			Id:                 k,
			Name:               v.Name,
			SumAmount:          sumAmount,
			AverageCoefficient: averageCoefficient,
			NumberOfBets:       int64(len(v.Coefficients)),
		}
		if s.BetAlreadyExist(ctx, k) {
			err = s.UpdateUserBetSummary(ctx, bet)
		} else {
			err = s.CreateUserBetSummary(ctx, bet)
		}
		if err != nil {
			return err
		}
	}
	for k, v := range blockedBetHistory {
		sumAmount := 0.0
		sumCoefficients := 0.0
		for _, value := range v.Bets {
			sumAmount += value
		}
		for _, value := range v.Coefficients {
			sumCoefficients += value
		}
		averageCoefficient := sumCoefficients / float64(len(v.Coefficients))
		bet := &blockedBets.BlockedBet{
			Id:                 k,
			Name:               v.Name,
			SumAmount:          sumAmount,
			AverageCoefficient: averageCoefficient,
			NumberOfBets:       int64(len(v.Coefficients)),
		}
		if s.BlockedBetAlreadyExist(ctx, k) {
			err = s.UpdateBlockedUserBetSummary(ctx, bet)
		} else {
			err = s.CreateBlockedUserBetSummary(ctx, bet)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) BetAlreadyExist(ctx context.Context, id int64) bool {
	bet, err := s.betsRps.Get(ctx, id)
	if bet != nil && err == nil {
		return true
	}
	return false
}
func (s *Service) BlockedBetAlreadyExist(ctx context.Context, id int64) bool {
	bet, err := s.blockedBetsRps.Get(ctx, id)
	if bet != nil && err == nil {
		return true
	}
	return false
}
func (s *Service) CreateUserBetSummary(ctx context.Context, bet *bets.Bet) error {
	return s.betsRps.Create(ctx, bet)
}
func (s *Service) CreateBlockedUserBetSummary(ctx context.Context, bet *blockedBets.BlockedBet) error {
	return s.blockedBetsRps.Create(ctx, bet)
}
func (s *Service) UpdateUserBetSummary(ctx context.Context, bet *bets.Bet) error {
	latestBet, err := s.betsRps.Get(ctx, bet.Id)
	if err != nil {
		return err
	}
	newNumberOfBets := bet.NumberOfBets + latestBet.NumberOfBets
	bet.AverageCoefficient = bet.AverageCoefficient*(float64(bet.NumberOfBets)/
		(float64(newNumberOfBets))) + latestBet.AverageCoefficient*(float64(latestBet.NumberOfBets)/
		(float64(newNumberOfBets)))
	bet.SumAmount = bet.SumAmount + latestBet.SumAmount
	bet.NumberOfBets = newNumberOfBets
	return s.betsRps.Update(ctx, bet)
}
func (s *Service) UpdateBlockedUserBetSummary(ctx context.Context, bet *blockedBets.BlockedBet) error {
	latestBet, err := s.blockedBetsRps.Get(ctx, bet.Id)
	if err != nil {
		return err
	}
	newNumberOfBets := bet.NumberOfBets + latestBet.NumberOfBets
	bet.AverageCoefficient = bet.AverageCoefficient*(float64(bet.NumberOfBets)/
		(float64(newNumberOfBets))) + latestBet.AverageCoefficient*(float64(latestBet.NumberOfBets)/
		(float64(newNumberOfBets)))
	bet.SumAmount = bet.SumAmount + latestBet.SumAmount
	bet.NumberOfBets = newNumberOfBets
	return s.blockedBetsRps.Update(ctx, bet)
}

//It is not necessary to create and return file
func (s *Service) Report(ctx context.Context) ([]byte, error) {
	contentOfFile := ""
	if s.processing == true {
		return nil, fmt.Errorf("Processing running ")
	}
	bets, err := s.betsRps.Listing(ctx)
	if err != nil {
		return nil, err
	}
	for i := range bets {
		str := fmt.Sprintf("%d, %s, %f, %f\n",
			bets[i].Id, bets[i].Name, bets[i].SumAmount, bets[i].AverageCoefficient)
		contentOfFile += str
	}
	if contentOfFile == "" {
		return nil, fmt.Errorf("Processing isn't start ")
	}
	return []byte(contentOfFile), nil
}
