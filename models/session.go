package models

import "time"

type Session struct {
    ID        int       `db:"id" json:"id"`
    NID       string    `db:"nid" json:"nid"`
    Token     string    `db:"token" json:"token"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
