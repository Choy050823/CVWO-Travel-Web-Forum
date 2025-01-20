// models/thread.go
package models

import "time"

type Comment struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	UserID     int       `json:"user_id"`
	ThreadID int       `json:"thread_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpVotes int `json:"upvotes"`
	DownVotes int `json:"downvotes"`
}