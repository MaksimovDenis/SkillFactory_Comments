package api

import (
	"net/http"
	"skillfactory/finalProject/commentsService/internal/api/oapi"
	"skillfactory/finalProject/commentsService/internal/storage"
	"skillfactory/finalProject/commentsService/internal/storage/queries"

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

	res := oapi.GetAllCommentsResponse{}

	if len(comments) != 0 {
		list := make([]oapi.Comment, len(comments))

		for idx, val := range comments {
			comment := oapi.Comment{
				Id:      int(val.ID),
				NewsId:  int(val.NewsID.Int32),
				Content: val.Content,
			}

			if val.ParentCommentID.Int32 != 0 {
				parentID := int(val.ParentCommentID.Int32)
				comment.ParentCommentId = &parentID
			}

			list[idx] = comment
		}

		res.JSON200 = &list
	}

	ctx.Header("Content-type", "application/json")

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No comments found"})
	}
}
