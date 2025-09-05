package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type FriendRequest struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Message   string    `json:"message"`
	Status    string    `json:"status"` // "pending", "accepted", "rejected"
	Timestamp time.Time `json:"timestamp"`
}

type FriendRepository struct {
	db *sqlx.DB
}

func NewFriendRepo(db *sqlx.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r FriendRepository) Get(from string, to string) (*FriendRequest, error) {
	profile := &FriendRequest{}
	err := r.db.QueryRowx("SELECT * FROM friend_request WHERE from = ? and to = ?", from, to).Scan(&profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

type FriendRequestProfile struct {
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	Status    string    `json:"status"` // "pending", "accepted", "rejected"
	Timestamp time.Time `json:"timestamp"`
}

func (r FriendRepository) ListTo(from string) (*[]FriendRequestProfile, error) {
	frs := []FriendRequestProfile{}
	err := r.db.Select(&frs, `SELECT
			p.name, fr.message, fr.status, fr.timestamp
		FROM friend_request fr
		JOIN profiles p ON fr.[to]=p.[id]
		WHERE [to] = ?`, from)
	if err != nil {
		return nil, err
	}
	return &frs, nil
}

func (r FriendRepository) SaveFriendRequest(fr FriendRequest) error {
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO friend_request (from, to, message, status, timestamp)
		VALUES (?, ?, ?, ?, ?, datetime('now'))`,
		fr)
	return err
}
