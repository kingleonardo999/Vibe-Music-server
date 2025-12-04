package repo

import (
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/db"
)

type CommentRepo struct{}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{}
}

func (c *CommentRepo) AddComment(comment *entity.Comment) error {
	return db.Get().Create(comment).Error
}

func (c *CommentRepo) GetCommentById(comment *entity.Comment, id uint64) error {
	return db.Get().First(comment, "id = ?", id).Error
}

func (c *CommentRepo) DeleteCommentById(id uint64) error {
	return db.Get().Delete(&entity.Comment{}, "id = ?", id).Error
}

func (c *CommentRepo) UpdateComment(comment *entity.Comment) error {
	return db.Get().Updates(comment).Error
}
