package handlers

import (
	"contatos/model"
	"contatos/templates"
	"contatos/util"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func ProfileEdit(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	profileRepo := model.NewProfileRepo(db)
	profile, err := profileRepo.Get(userId)
	if err != nil {
		profile = &model.Profile{}
	}

	templates.Profile(*profile).Render(r.Context(), w)
}

func ProfileSave(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Create profile with form data
	profile := model.Profile{
		ID:   userId,
		Name: r.FormValue("first_name"),
		Data: make(map[string]string),
	}

	// Handle custom key-value pairs
	customKeys := r.Form["custom_key"]
	customValues := r.Form["custom_value"]
	for i, key := range customKeys {
		if i < len(customValues) && key != "" {
			profile.Data[key] = customValues[i]
		}
	}

	profileRepo := model.NewProfileRepo(db)
	profileRepo.Save(profile)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ProfileFind(db *sqlx.DB, w http.ResponseWriter, r *http.Request, config util.Config, userId string) {

	profileRepo := model.NewProfileRepo(db)
	profiles, _ := profileRepo.Search(r.FormValue("q"))
	// if err != nil {
	// 	profiles = &model.Profile{}
	// }

	templates.ProfileFindResult(*profiles).Render(r.Context(), w)
}
