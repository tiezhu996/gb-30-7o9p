package model

import "time"

// Favorite links a user to a pet or a post.
type Favorite struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index:idx_fav_user_target,unique;not null" json:"user_id"`
	TargetType string    `gorm:"size:16;index:idx_fav_user_target,unique;not null" json:"target_type"`
	TargetID   uint      `gorm:"index:idx_fav_user_target,unique;not null" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}
