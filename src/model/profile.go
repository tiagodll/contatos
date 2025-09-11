package model

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type Profile struct {
	ID   string            `json:"id"`
	Name string            `json:"name"`
	Data map[string]string `json:"ignore"`
}

type dbProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewProfileRepo(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

type ProfileRepository struct {
	db *sqlx.DB
}

func (r ProfileRepository) Get(id string) (*Profile, error) {
	profile := &Profile{}
	var data string
	err := r.db.QueryRow("SELECT id, name, data FROM profiles WHERE id = ?", id).Scan(&profile.ID, &profile.Name, &data)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(data), &profile.Data)
	return profile, nil
}

func (r ProfileRepository) Save(p Profile) error {
	data, _ := json.Marshal(p.Data)
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO profiles (id, name, data, timestamp)
		VALUES (?, ?, ?, datetime('now'))`,
		p.ID, p.Name, data)
	return err
}

func (r ProfileRepository) Search(q string, userId string) (*[]Profile, error) {
	profiles := []Profile{}
	rows, err := r.db.Queryx("SELECT id, name, data FROM profiles WHERE id != ? AND name LIKE ?", userId, q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p dbProfile
		_ = rows.StructScan(&p)
		profile := &Profile{
			ID:   p.ID,
			Name: p.Name,
		}
		json.Unmarshal([]byte(p.Data), &profile.Data)
		profiles = append(profiles, *profile)
	}
	return &profiles, nil
}
