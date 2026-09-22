package dto

// FavoriteRequest adds a favorite.
type FavoriteRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   uint   `json:"target_id" binding:"required"`
}
