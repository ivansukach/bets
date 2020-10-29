package blocking

import (
	"errors"
	"net/http"
)

type BlockUsersReqModel struct {
	Ids []int64 `json:"ids"`
}

func (b *BlockUsersReqModel) Bind(r *http.Request) error {
	if len(b.Ids) == 0 {
		return errors.New("Field Ids should contain at least one user ")
	}
	return nil
}
