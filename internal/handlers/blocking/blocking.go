package blocking

import "github.com/ivansukach/bets/internal/service/blocking"

type Blocking struct {
	blockingService *blocking.Service
}

func New(blockingService *blocking.Service) *Blocking {
	return &Blocking{blockingService: blockingService}
}
