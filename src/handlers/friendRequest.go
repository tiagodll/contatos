package handlers

import (
	"contatos/model"
	"contatos/templates"
	"contatos/util"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func FriendRequest(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	profileId := r.URL.Query().Get("profileId")
	fr := &model.FriendRequest{
		From:      userId,
		To:        profileId,
		Status:    "requested",
		Timestamp: time.Now(),
	}

	profileRepo := model.NewProfileRepo(db)
	p, err := profileRepo.Get(profileId)
	if err != nil {
		log.Printf("Invalid form data: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	templates.FriendRequest(fr, p).Render(r.Context(), w)
}

func FriendRequestSave(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	err := r.ParseForm()
	if err != nil {
		log.Printf("Invalid form data: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Create profile with form data
	fr := model.FriendRequest{
		From:      userId,
		To:        r.FormValue("to"),
		Message:   r.FormValue("message"),
		Status:    "requested",
		Timestamp: time.Now(),
	}

	fRepo := model.NewFriendRepo(db)
	err = fRepo.SaveFriendRequest(fr)
	if err != nil {
		log.Printf("Saving friend request: %v", err)
		http.Error(w, "Saving friend request error", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AcceptFriendRequest(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	fr := model.FriendRequest{
		From: r.PathValue("from"),
		To:   userId,
	}

	fRepo := model.NewFriendRepo(db)
	err := fRepo.AcceptFriendRequest(fr)
	if err != nil {
		log.Printf("Accepting friend request: %v", err)
		http.Error(w, "Accepting friend request error", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RejectFriendRequest(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	fr := model.FriendRequest{
		From: r.PathValue("from"),
		To:   userId,
	}

	fRepo := model.NewFriendRepo(db)
	err := fRepo.RejectFriendRequest(fr)
	if err != nil {
		log.Printf("Accepting friend request: %v", err)
		http.Error(w, "Accepting friend request error", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
