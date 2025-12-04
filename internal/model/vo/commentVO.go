package vo

import "time"

type CommentVO struct {
	CommentID  uint64    `json:"commentId"`
	Username   string    `json:"username"`
	UserAvatar string    `json:"userAvatar"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createTime" time_format:"2006-01-02"` // 仅日期
	LikeCount  uint64    `json:"likeCount"`
}
