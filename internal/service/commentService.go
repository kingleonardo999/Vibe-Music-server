package service

import (
	"time"
	"vibe-music-server/internal/model/dto"
	"vibe-music-server/internal/model/entity"
	"vibe-music-server/internal/pkg/result"
	"vibe-music-server/internal/pkg/result/consts"
	"vibe-music-server/internal/pkg/util"
	"vibe-music-server/internal/repo"
)

type CommentService struct {
	commentRepo *repo.CommentRepo
}

func NewCommentService(commentRepo *repo.CommentRepo) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (c CommentService) AddSongComment(commentSongDTO *dto.CommentSongDTO, claims *util.Claims) result.Result[result.Nil] {
	retErr := result.Error[result.Nil]
	retSuc := result.Success[result.Nil]
	userID := claims.UserId
	comment := entity.Comment{
		UserID:     userID,
		SongID:     &commentSongDTO.SongID,
		Content:    commentSongDTO.Content,
		CreateTime: time.Now(),
		Type:       entity.CommentTypeSong,
		LikeCount:  0,
	}
	if err := c.commentRepo.AddComment(&comment); err != nil {
		return retErr(consts.Add + consts.Failed)
	}
	util.DeleteCacheByPattern("song:*")
	return retSuc(consts.Add + consts.Success)
}

func (c CommentService) AddPlaylistComment(commentPlaylistDTO *dto.CommentPlaylistDTO, claims *util.Claims) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	userID := claims.UserId
	comment := entity.Comment{
		UserID:     userID,
		PlaylistID: &commentPlaylistDTO.PlaylistID,
		Content:    commentPlaylistDTO.Content,
		CreateTime: time.Now(),
		Type:       entity.CommentTypePlaylist,
		LikeCount:  0,
	}
	if err := c.commentRepo.AddComment(&comment); err != nil {
		return retErr(consts.Add + consts.Failed)
	}
	util.DeleteCacheByPattern("playlist:*")
	return retSuc(consts.Add + consts.Success)
}

func (c CommentService) LikeComment(commentId uint64) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var comment entity.Comment
	if err := c.commentRepo.GetCommentById(&comment, commentId); err != nil {
		return retErr(consts.DataNotFound)
	}
	comment.LikeCount++
	if err := c.commentRepo.UpdateComment(&comment); err != nil {
		return retErr(consts.Update + consts.Failed)
	}
	if comment.Type == entity.CommentTypeSong {
		util.DeleteCacheByPattern("song:*")
	} else if comment.Type == entity.CommentTypePlaylist {
		util.DeleteCacheByPattern("playlist:*")
	}
	return retSuc(consts.Update + consts.Success)
}

func (c CommentService) CancelLikeComment(commentId uint64) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	var comment entity.Comment
	if err := c.commentRepo.GetCommentById(&comment, commentId); err != nil {
		return retErr(consts.DataNotFound)
	}
	if comment.LikeCount > 0 {
		comment.LikeCount--
		if err := c.commentRepo.UpdateComment(&comment); err != nil {
			return retErr(consts.Update + consts.Failed)
		}
		if comment.Type == entity.CommentTypeSong {
			util.DeleteCacheByPattern("song:*")
		} else if comment.Type == entity.CommentTypePlaylist {
			util.DeleteCacheByPattern("playlist:*")
		}
	}
	return retSuc(consts.Update + consts.Success)
}

func (c CommentService) DeleteComment(commentId uint64, claims *util.Claims) result.Result[result.Nil] {
	var retErr = result.Error[result.Nil]
	var retSuc = result.Success[result.Nil]
	userID := claims.UserId
	var comment entity.Comment
	if err := c.commentRepo.GetCommentById(&comment, commentId); err != nil {
		return retErr(consts.DataNotFound)
	}
	if comment.UserID != userID {
		return retErr(consts.NoPermission)
	}
	if err := c.commentRepo.DeleteCommentById(commentId); err != nil {
		return retErr(consts.Delete + consts.Failed)
	}
	if comment.Type == entity.CommentTypeSong {
		util.DeleteCacheByPattern("song:*")
	} else if comment.Type == entity.CommentTypePlaylist {
		util.DeleteCacheByPattern("playlist:*")
	}
	return retSuc(consts.Delete + consts.Success)
}
