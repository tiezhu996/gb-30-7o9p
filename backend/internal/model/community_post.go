package model

import "time"

// CommunityPost is a rescue story or lost-pet notice.
type CommunityPost struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	OrgID        uint      `gorm:"index" json:"org_id"`
	Title        string    `gorm:"size:255;not null" json:"title"`
	Content      string    `gorm:"type:text" json:"content"`
	Images       string    `gorm:"type:json" json:"images"`
	PostType     string    `gorm:"size:16;default:story;index" json:"post_type"`
	LikeCount    int       `gorm:"default:0" json:"like_count"`
	CommentCount int       `gorm:"default:0" json:"comment_count"`
	Status       string    `gorm:"size:16;default:published" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
