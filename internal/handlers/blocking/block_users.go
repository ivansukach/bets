package blocking

import (
	"context"
	"github.com/go-chi/render"
	"github.com/ivansukach/bets/internal/tools"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (b *Blocking) BlockUsers(w http.ResponseWriter, r *http.Request) {
	users := &BlockUsersReqModel{}
	if err := render.Bind(r, users); err != nil {
		log.Error(err)
		render.Render(w, r, tools.ErrInvalidRequest(err))
		return
	}
	log.Debugf("users to block: %v", *users)
	err := b.blockingService.BlockUsers(context.Background(), users.Ids)
	if err != nil {
		log.Error(err)
		render.Render(w, r, tools.InternalServerError(err))
		return
	}

	render.Status(r, http.StatusOK)
}
