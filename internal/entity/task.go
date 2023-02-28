package entity

import "time"

type Link struct {
	Hash      string `json:"hash"`
	ToLink    string `json:"to_link"`
	CreatedAt time.Time
}
