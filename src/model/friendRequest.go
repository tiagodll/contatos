package model

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type FriendRequestReository struct {
	db *sqlx.DB
}

func NewFriendRequestRepo(db *sqlx.DB) *FriendRequestReository {
	return &FriendRequestReository{db: db}
}

type FriendRequest struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Message   string    `json:"message"`
	Status    string    `json:"status"` // "pending", "accepted", "rejected"
	Timestamp time.Time `json:"timestamp"`
}

func (r FriendRepository) GetRequest(from string, to string) (*FriendRequest, error) {
	profile := &FriendRequest{}
	err := r.db.QueryRowx("SELECT * FROM friend_requests WHERE from = ? and to = ?", from, to).Scan(&profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

type FriendRequestProfile struct {
	Name      string `json:"name"`
	From      string `json:"from"`
	Message   string `json:"message"`
	Status    string `json:"status"` // "pending", "accepted", "rejected"
	Timestamp string `json:"timestamp"`
}

func (r FriendRepository) ListFriendRequests(to string) (*[]FriendRequestProfile, error) {
	frs := []FriendRequestProfile{}
	err := r.db.Select(&frs, `SELECT
			p.name, fr.[from], fr.message, fr.status, fr.timestamp
		FROM friend_requests fr
		JOIN profiles p ON fr.[to]=p.[id]
		WHERE [to] = ?`, to)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
		return nil, err
	}
	return &frs, nil
}
func (r FriendRepository) SaveFriendRequest(fr FriendRequest) error {
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO friend_requests ([from], [to], [message], [status], [timestamp])
		VALUES (?, ?, ?, ?, datetime('now'))`,
		fr.From, fr.To, fr.Message, fr.Status)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}
	return err
}

func (r FriendRepository) AcceptFriendRequest(fr FriendRequest) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO friends ([user_id], [friend_id], [notes]) VALUES (?, ?, '')`, fr.From, fr.To)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	_, err = r.db.Exec(`INSERT OR REPLACE INTO friends ([user_id], [friend_id], [notes]) VALUES (?, ?, '')`, fr.To, fr.From)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	_, err = r.db.Exec(`DELETE FROM friend_requests WHERE [from]=? AND [to]=?`, fr.From, fr.To)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	return err
}

func (r FriendRepository) RejectFriendRequest(fr FriendRequest) error {
	_, err := r.db.Exec(`DELETE FROM friend_requests WHERE [from]=? AND [to]=?`, fr.From, fr.To)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	return err
}
