package handlers

import (
	"contatos/model"
	"contatos/templates"
	"contatos/util"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func IndexHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config) {

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		templates.LandingPage(config.GetLoginUrl()).Render(r.Context(), w)
		return
	}

	token := cookie.Value
	tokendata, err := util.DecodeJWT(token, config.Jwt.Secret)
	if err != nil {
		templates.LandingPage(config.GetLoginUrl()).Render(r.Context(), w)
		return
	}

	user_id := tokendata["sub"].(string)

	profileRepo := model.NewProfileRepo(db)
	profile, err := profileRepo.Get(user_id)
	if err != nil {
		// user should always have a profile, if not, create it
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	friendsRepo := model.NewFriendRepo(db)
	requests, err := friendsRepo.ListTo(user_id)

	templates.Dashboard(*profile, *requests, config.GetLoginUrl()).Render(r.Context(), w)
}
