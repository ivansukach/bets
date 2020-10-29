package bets

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/ivansukach/bets/internal/tools"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"time"
)

func (b *Bets) Report(w http.ResponseWriter, r *http.Request) {
	fileContent, err := b.betsService.Report(context.Background())
	if err != nil {
		if err.Error() == "Processing running " || err.Error() == "Processing isn't start " {
			render.Render(w, r, &tools.MessageContainer{
				Message: err.Error(),
			})
			return
		}
		log.Error(err)
		render.Render(w, r, tools.InternalServerError(err))
		return
	}

	//render.Status(r, http.StatusOK)
	//f:=&tools.MessageContainer{File: file}

	writer := multipart.NewWriter(w)
	w.Header().Set("Content-type", writer.FormDataContentType())
	part, err := writer.CreateFormFile("file", fmt.Sprintf("%d%s", time.Now().UTC().UnixNano(), ".csv"))
	if err != nil {
		log.Error(err)
		render.Render(w, r, tools.InternalServerError(err))
		return
	}
	part.Write(fileContent)
}
