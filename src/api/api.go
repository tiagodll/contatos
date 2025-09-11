package api

import (
	"contatos/model"
	"contatos/util"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func GetContactList(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	repo := model.NewFriendRepo(db)
	friends, err := repo.ListFriendNames(userId)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	util.RespondWithSuccess(w, friends)
}
