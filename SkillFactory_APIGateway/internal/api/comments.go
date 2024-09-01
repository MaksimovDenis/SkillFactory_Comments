package api

import (
	"encoding/json"
	"io"
	"net/http"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api/oapi"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	commentsServiceURL = "http://localhost:8882"
)

func (api *API) CreateComment(ctx *gin.Context) {
	rout := "/api/comments"

	resp, err := http.Post(commentsServiceURL+rout, "application/json", ctx.Request.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get comments from comments service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments from comments service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get comments from comments service"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	ctx.Data(resp.StatusCode, "application/json", body)
}

func (api *API) GetAllComments(ctx *gin.Context) {

	rout := "/api/comments"

	resp, err := http.Get(commentsServiceURL + rout)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get comments from comments service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments from comments service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get comments from comments service"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	res := oapi.GetAllCommentsResponse{}

	var comments []models.Comments

	err = json.Unmarshal(body, &comments)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to unmarshal response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response body"})
		return
	}

	list := make([]oapi.Comment, len(comments))

	for idx, val := range comments {
		comment := oapi.Comment{
			Id:              int(val.ID),
			NewsId:          int(val.NewsID),
			ParentCommentId: val.ParentCommentID,
			Content:         val.Content,
		}

		list[idx] = comment
	}

	res.JSON200 = &list

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No comments found"})
	}
}
