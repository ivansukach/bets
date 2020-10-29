package blocking

import (
	"context"
	"fmt"
	users "github.com/ivansukach/bets/internal/repositories/blocked-users"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	blockedUsersRps users.Repository
}

func New(blockedUsersRps users.Repository) *Service {
	return &Service{blockedUsersRps: blockedUsersRps}
}
func (s *Service) BlockUsers(ctx context.Context, ids []int64) error {
	for i := range ids {
		if err := s.blockedUsersRps.Create(ctx, &users.BlockedUser{Id: ids[i]}); err != nil {
			log.Error(err)
		}
	}
	return nil
}
func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	return s.blockedUsersRps.Delete(ctx, id)
}
func (s *Service) BlockedReport(ctx context.Context) ([]byte, error) {
	users, err := s.blockedUsersRps.Listing(ctx)
	if err != nil {
		return nil, err
	}
	contentOfFile := ""
	for i := range users {
		str := fmt.Sprintf("%d\n", users[i].Id)
		contentOfFile += str
	}
	return []byte(contentOfFile), nil
}
