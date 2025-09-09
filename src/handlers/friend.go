package handlers

import (
	"contatos/model"
	"contatos/util"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func UpdateNotes(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	err := r.ParseForm()
	if err != nil {
		log.Printf("Invalid form data: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	fr := model.Friend{
		FriendId: r.PathValue("id"),
		UserId:   userId,
		Notes:    r.FormValue("notes"),
	}

	repo := model.NewFriendRepo(db)
	err = repo.SaveFriend(fr)
	if err != nil {
		log.Printf("Accepting friend request: %v", err)
		http.Error(w, "Accepting friend request error", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
