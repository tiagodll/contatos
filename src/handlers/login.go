package handlers

import (
	"contatos/util"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func LoginHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request, c util.Config) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",                              // Cookie available site-wide
		Expires:  time.Now().Add(-1 * time.Second), // Match JWT expiration
		HttpOnly: false,                            // Prevent JavaScript access
		Secure:   false,                            // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,             // Protect against CSRF
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, c.GetLoginUrl(), http.StatusFound)
}

func LoginCallbackHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request, c util.Config) {

	token := r.URL.Query().Get("token")

	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",                            // Cookie available site-wide
		Expires:  time.Now().Add(24 * time.Hour), // Match JWT expiration
		HttpOnly: false,                          // Prevent JavaScript access
		Secure:   false,                          // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,           // Protect against CSRF
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// // generateTokenHandler generates a JWT token for testing purposes
// func generateTokenHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		util.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
// 		return
// 	}

// 	var requestBody struct {
// 		UserID string `json:"user_id"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
// 		util.RespondWithError(w, http.StatusBadRequest, "Invalid JSON body")
// 		return
// 	}

// 	if requestBody.UserID == "" {
// 		util.RespondWithError(w, http.StatusBadRequest, "user_id is required")
// 		return
// 	}

// 	token, err := util.GenerateJWT(requestBody.UserID, config.Jwt)
// 	if err != nil {
// 		log.Printf("Error generating JWT: %v", err)
// 		util.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
// 		return
// 	}

// 	util.RespondWithSuccess(w, map[string]interface{}{
// 		"token":   token,
// 		"user_id": requestBody.UserID,
// 	})
// }
