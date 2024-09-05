package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api/oapi"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	feedsServiceURL = "http://localhost:8883/api/feeds"
)

func (api *API) Feeds(ctx *gin.Context, params oapi.FeedsParams) {
	parsedURL, err := url.Parse(feedsServiceURL)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to parse url address feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse url address feeds service"})
		return
	}

	if *params.Page <= 0 {
		*params.Page = 1
	}

	paramsQuery := url.Values{}
	paramsQuery.Add("page", strconv.FormatInt(int64(*params.Page), 10))
	paramsQuery.Add("title", *params.Title)
	paramsQuery.Add("filter", *params.Filter)

	parsedURL.RawQuery = paramsQuery.Encode()

	reqURL := parsedURL.String()

	fmt.Println(reqURL)
	resp, err := http.Get(reqURL)
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

	resp, err := http.Get(feedsServiceURL + "/" + idStr)
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
