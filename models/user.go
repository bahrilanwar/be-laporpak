package models

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	NID       string    `db:"nid" json:"nid"`
	Password  string    `db:"password" json:"password"`
	NamaUser  string    `db:"nama_user" json:"nama_user"`
	Roles     string    `db:"roles" json:"roles"`
	Level     string    `db:"level" json:"level"`
	Team      string    `db:"team" json:"team"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
