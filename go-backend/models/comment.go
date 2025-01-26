// models/comment.go
package models

import "time"

type Comment struct {
	ID             int       `json:"id"`
	Content        string    `json:"content"`
	UserID         int       `json:"userId"`
	ThreadID       int       `json:"threadId"`
	AttachedImages []string  `json:"attachedImages"`
	Upvotes        int       `json:"upvotes"`
	Downvotes      int       `json:"downvotes"`
	CreatedAt      time.Time `json:"createdAt"`
	Author         string    `json:"author"`
}
