package api

import (
	"net/http"

	"github.com/MaksimovDenis/SkillFactory_Comments/internal/storage"
	"github.com/MaksimovDenis/SkillFactory_Comments/internal/storage/queries"

	"github.com/gin-gonic/gin"
)

func (api *API) CreateComment(ctx *gin.Context) {
	var commentCreate storage.Comment

	if err := ctx.BindJSON(&commentCreate); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshal comment body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse comment object body"})

		return
	}

	args := queries.CreateCommentsParams{
		NewsID:          commentCreate.NewsID,
		ParentCommentID: commentCreate.ParentCommentID,
		Content:         commentCreate.Content,
	}

	if err := api.storage.Queries.CreateComments(ctx.Request.Context(), args); err != nil {
		api.l.Error().Err(err).Msg("failed to add comment to storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment to storage"})
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (api *API) GetAllComments(ctx *gin.Context) {
	comments, err := api.storage.Queries.GetAllComments(ctx.Request.Context())
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get comments from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments from storage"})

		return
	}

	dataList := make([]storage.Comment, len(comments))

	if len(comments) != 0 {

		for idx, val := range comments {
			comment := storage.Comment{
				ID:              val.ID,
				NewsID:          val.NewsID,
				ParentCommentID: val.ParentCommentID,
				Content:         val.Content,
			}

			dataList[idx] = comment
		}

		ctx.JSON(http.StatusOK, dataList)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No comments found"})
	}
}
