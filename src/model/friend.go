package model

import (
	"database/sql"
)

type Friend struct {
	ID   string            `json:"id"`
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

type FriendRepository struct {
	db *sql.DB
}

func (r FriendRepository) Get(id string) (*Friend, error) {
	profile := &Friend{}
	err := r.db.QueryRow("SELECT id, name, email, phone, data FROM friends WHERE id = ?", id).Scan(&profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r FriendRepository) SaveFriend(profile Friend) error {
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO friends (id, name, email, phone, data, timestamp)
		VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		profile)
	return err
}
