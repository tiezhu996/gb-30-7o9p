package dto

// PostCreateRequest creates a community post.
type PostCreateRequest struct {
	Title    string `json:"title" binding:"required,max=255"`
	Content  string `json:"content" binding:"required"`
	Images   string `json:"images"`
	PostType string `json:"post_type" binding:"omitempty,oneof=story lost_notice"`
}

// CommentCreateRequest adds a comment.
type CommentCreateRequest struct {
	Content string `json:"content" binding:"required,max=1000"`
}
