package api

import (
	"encoding/json"
	"io"
	"net/http"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api/oapi"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	feedsServiceURL = "http://localhost:8883"
)

func (api *API) Feeds(ctx *gin.Context, params oapi.FeedsParams) {
	if *params.Limit <= 0 {
		*params.Limit = 10
	}

	limitStr := strconv.FormatInt(int64(*params.Limit), 10)

	rout := "/api/feeds/" + limitStr

	resp, err := http.Get(feedsServiceURL + rout)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds from feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds from feeds service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get feeds from feeds service"})
		return
	}

	res, err := oapi.ParseFeedsResponse(resp)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No feeds found"})
	}
}

func (api *API) FeedsById(ctx *gin.Context, id oapi.ID) {
	if id <= 0 {
		api.l.Error().Msg("Invalid ID argument")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID argument"})
		return
	}

	idStr := strconv.FormatInt(int64(id), 10)

	rout := "/api/feeds/" + idStr

	resp, err := http.Get(feedsServiceURL + rout)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds by id from feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds by id from feeds service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get feeds by id from feeds service"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var feed oapi.Feeds
	var res oapi.FeedsByIdResponse

	err = json.Unmarshal(body, &feed)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to unmarshal response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response body"})
		return
	}

	list := make([]oapi.Feeds, 1)
	list[0] = feed

	res.JSON200 = &list

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No feeds found"})
	}
}

func (api *API) FeedsByFilter(ctx *gin.Context) {
	rout := "/api/feeds"

	resp, err := http.Post(feedsServiceURL+rout, "application/json", ctx.Request.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds by filter from feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds by filter from feeds service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "FFailed to get feeds by filter from feeds service"})
		return
	}

	res, err := oapi.ParseFeedsByIdResponse(resp)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No feeds found"})
	}
}
