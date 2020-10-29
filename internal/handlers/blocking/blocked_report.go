package blocking

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

func (b *Blocking) BlockedReport(w http.ResponseWriter, r *http.Request) {
	fileContent, err := b.blockingService.BlockedReport(context.Background())
	if err != nil {
		log.Error(err)
		render.Render(w, r, tools.InternalServerError(err))
		return
	}

	writer := multipart.NewWriter(w)
	defer writer.Close()
	w.Header().Set("Content-type", writer.FormDataContentType())
	part, err := writer.CreateFormFile("file", fmt.Sprintf("%d%s", time.Now().UTC().UnixNano(), ".csv"))
	if err != nil {
		log.Error(err)
		render.Render(w, r, tools.InternalServerError(err))
		return
	}
	part.Write(fileContent)
}
