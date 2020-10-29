package bets

import (
	"context"
	"github.com/go-chi/render"
	"github.com/ivansukach/bets/internal/tools"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (b *Bets) Process(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Error(err)
		render.Render(w, r, tools.ErrInvalidRequest(err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error(err)
		render.Render(w, r, tools.ErrInvalidRequest(err))
		return
	}
	defer file.Close()
	log.Debugf("Uploaded File: %+v\n", header.Filename)
	log.Debugf("File Size: %+v\n", header.Size)
	log.Debugf("MIME Header: %+v\n", header.Header)

	err = b.betsService.Process(context.Background(), &file)

	render.Status(r, http.StatusOK)
}
