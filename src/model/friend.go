package model

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
)

type Friend struct {
	UserId   string `json:"user_id"`
	FriendId string `json:"friend_id"`
	Notes    string `json:"notes"`
}
type FriendName struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Image *string `json:"image"`
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

func (r FriendRepository) GetFriend(userId string, friendId string) (*FriendProfile, error) {
	p := dbFriendProfile{}
	err := r.db.QueryRowx(`SELECT
			p.id, p.name, p.data, f.notes
		FROM friends f
			JOIN profiles p ON f.[friend_id]=p.[id]
		WHERE f.[user_id] = ? and f.[friend_id] = ?`, userId, friendId).StructScan(&p)
	if err != nil {
		return nil, err
	}
	profile := &FriendProfile{
		ID:    p.ID,
		Name:  p.Name,
		Notes: p.Notes,
	}
	json.Unmarshal([]byte(p.Data), &profile.Data)

	return profile, nil
}

func (r FriendRepository) SaveFriend(fr Friend) error {
	_, err := r.db.Exec(`UPDATE friends SET notes = ? WHERE user_id=? AND friend_id=?`, fr.Notes, fr.UserId, fr.FriendId)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}
	return err
}

func (r FriendRepository) ListFriends(userId string) (*[]FriendProfile, error) {
	dbprofiles := []dbFriendProfile{}
	err := r.db.Select(&dbprofiles, `SELECT
			p.id, p.name, p.data, f.notes
		FROM friends f
		JOIN profiles p ON f.[friend_id]=p.[id]
		WHERE f.[user_id] = ?`, userId)
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
func (r FriendRepository) Unfriend(fr Friend) error {
	_, err := r.db.Exec(`DELETE FROM friends WHERE [user_id]=? AND [friend_id] = ?`, fr.UserId, fr.FriendId)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	_, err = r.db.Exec(`DELETE FROM friends WHERE [user_id]=? AND [friend_id] = ?`, fr.FriendId, fr.UserId)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
	}

	return err
}

func (r FriendRepository) ListFriendNames(userId string) (*[]FriendName, error) {
	list := []FriendName{}
	err := r.db.Select(&list,
		`SELECT p.id, p.name, p.image
		FROM friends f
		JOIN profiles p ON f.[friend_id]=p.[id]
		WHERE f.[user_id] = ?`, userId)
	if err != nil {
		log.Printf("Error querying the db: %v", err)
		return nil, err
	}
	return &list, nil
}
