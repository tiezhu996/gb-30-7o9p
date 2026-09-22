package model

import "time"

// PostComment is a comment under a community post.
type PostComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"index;not null" json:"post_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Content   string    `gorm:"size:1000;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
