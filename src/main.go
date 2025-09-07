package main

import (
	"contatos/handlers"
	"contatos/util"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := util.GetEnvConfig()

	var err error
	db, err := sqlx.Open("sqlite3", config.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	util.CreateTables(db)

	Handle := func(url string, handler func(db *sqlx.DB, w http.ResponseWriter, r *http.Request, c util.Config)) {
		http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
			handler(db, w, r, config)
		})
	}

	HandleAuth := func(url string, handler func(db *sqlx.DB, w http.ResponseWriter, r *http.Request, c util.Config, userID string)) {
		http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("session")
			if err == http.ErrNoCookie {
				http.Error(w, "Cookie not found", http.StatusUnauthorized)
				return
			} else if err != nil {
				http.Error(w, "Error retrieving cookie", http.StatusBadRequest)
				return
			}

			// Basic validation: check if the cookie value is not empty
			if cookie.Value == "" {
				http.Error(w, "Empty cookie value", http.StatusUnauthorized)
				return
			}

			// Example: Validate cookie expiration (if stored in the cookie value or metadata)
			/*if cookie.Expires.Before(time.Now()) {
				http.Error(w, "Cookie expired", http.StatusUnauthorized)
				return
			}//*/

			userID, err := util.ValidateJWT(cookie.Value, config.Jwt)
			if err != nil {
				util.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
				return
			}

			r.Header.Set("X-User-ID", userID)
			handler(db, w, r, config, userID)
		})
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticPath))))

	Handle("/callback", handlers.LoginCallbackHandler)
	Handle("/login", handlers.LoginHandler)

	HandleAuth("GET /profile", handlers.ProfileEdit)
	HandleAuth("POST /profile", handlers.ProfileSave)
	HandleAuth("POST /find-profile", handlers.ProfileFind)
	HandleAuth("GET /friend-request", handlers.FriendRequest)
	HandleAuth("POST /friend-request", handlers.FriendRequestSave)

	// Auth endpoint for testing
	// Handle("/api/generate-token", func(w http.ResponseWriter, r *http.Request) { generateTokenHandler(db, w, r) })

	Handle("/", handlers.IndexHandler)

	currentpwd, _ := os.Getwd()
	fmt.Printf("Contacts App running from %s\n", currentpwd)
	fmt.Printf("Logs at %s\n", config.LogPath)
	fmt.Printf("DB at %s\n", config.DbConnectionString)
	fmt.Println("Server starting on " + config.AppPort)
	log.Fatal(http.ListenAndServe(config.AppPort, nil))
}
