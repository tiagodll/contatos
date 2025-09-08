package model

import (
	"encoding/json"
	"log"
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

type Friend struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Notes string `json:"notes"`
}
type FriendProfile struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Data  map[string]string `json:"data"`
	Notes string            `json:"notes"`
}
type dbFriendProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Data  string `json:"data"`
	Notes string `json:"notes"`
}

type FriendRepository struct {
	db *sqlx.DB
}

func NewFriendRepo(db *sqlx.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r FriendRepository) Get(from string, to string) (*FriendRequest, error) {
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

func (r FriendRepository) ListFriends(to string) (*[]FriendProfile, error) {
	dbprofiles := []dbFriendProfile{}
	err := r.db.Select(&dbprofiles, `SELECT
			p.id, p.name, p.data, f.notes
		FROM friends f
		JOIN profiles p ON f.[to]=p.[id]
		WHERE f.[to] = ?`, to)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
		return nil, err
	}

	profiles := []FriendProfile{}
	for _, p := range dbprofiles {
		profile := &FriendProfile{
			ID:    p.ID,
			Name:  p.Name,
			Notes: p.Notes,
		}
		json.Unmarshal([]byte(p.Data), &profile.Data)
		profiles = append(profiles, *profile)
	}
	return &profiles, nil
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
	_, err := r.db.Exec(`INSERT OR REPLACE INTO friends ([from], [to], [notes]) VALUES (?, ?, '')`, fr.From, fr.To)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	_, err = r.db.Exec(`INSERT OR REPLACE INTO friends ([from], [to], [notes]) VALUES (?, ?, '')`, fr.To, fr.From)
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
