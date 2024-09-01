package api

import (
	"net/http"
	"strconv"

	"github.com/MaksimovDenis/SkillFactory_News/internal/models"
	"github.com/gin-gonic/gin"
)

func (api *API) Feeds(ctx *gin.Context) {
	limitStr := ctx.Param("limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get news from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid limit parameter"})
		return
	}

	if limit == 0 || limit < 0 {
		limit = 10
	}

	feeds, err := api.storage.Feeds.Feeds(limit)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get news from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news from storage"})
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.JSON(http.StatusOK, feeds)
}

func (api *API) FeedById(ctx *gin.Context) {
	queryId := ctx.Param("id")

	id, err := strconv.Atoi(queryId)
	if err != nil {
		api.l.Error().Err(err).Msg("invalid arguments")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid arguments"})
		return
	}

	feed, err := api.storage.Feeds.FeedById(id)
	if err != nil {
		api.l.Error().Err(err).Msgf("failed to get news from storage:%v", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get news from storage"})
		return
	}

	ctx.JSON(http.StatusOK, feed)
}

func (api *API) FeedsByFilter(ctx *gin.Context) {
	var filters models.Filter

	if err := ctx.BindJSON(&filters); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshal feeds body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse feeds object body"})

		return
	}

	if filters.Limit == 0 || filters.Limit < 0 {
		filters.Limit = 10
	}

	feeds, err := api.storage.Feeds.FeedsByFilter(filters.Limit, filters.Filter)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds from storage"})
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.JSON(http.StatusOK, feeds)
}
