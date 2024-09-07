package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	period             = 5 // seconds
	deepPavlovModelURL = "http://0.0.0.0:5555/model"
)

type Censor struct {
	Content []string `json:"x"`
}

type CensorResult struct {
	Data []byte
	IDs  []int32
}

func (api *API) StartCensorship() {
	chResult := make(chan CensorResult)
	chErrors := make(chan error)

	go api.handleResults(chResult)
	go api.handleErrors(chErrors)
	go api.parseURL(chResult, chErrors)
}

func (api *API) handleResults(chResult <-chan CensorResult) {
	for result := range chResult {
		var censorResult [][]string
		err := json.Unmarshal(result.Data, &censorResult)
		if err != nil {
			api.l.Debug().Err(err).Msgf("Error unmarshalling JSON: %v", err)
			continue
		}

		for idx, res := range censorResult {
			if len(res) > 0 && res[0] == "negative" {
				err := api.storage.Queries.DeleteComment(context.Background(), result.IDs[idx])
				if err != nil {
					api.l.Err(err).Msgf("Failed to delete comment with ID %d: %v", result.IDs[idx], err)
				}
				api.l.Debug().Msgf("Comment with ID %d has been removed", result.IDs[idx])
			}
		}
	}
}

func (api *API) handleErrors(chErrors <-chan error) {
	for err := range chErrors {
		api.l.Debug().Err(err).Msg("Error processing request")
	}
}

func (api *API) getDataFromStorage() ([]byte, []int32, error) {
	comments, err := api.storage.Queries.GetAllComments(context.Background())
	if err != nil {
		api.l.Debug().Err(err).Msg("Failed to get comments from storage")
		return nil, nil, err
	}

	commentsData := make([]string, len(comments))
	idsData := make([]int32, len(comments))

	for idx, value := range comments {
		commentsData[idx] = value.Content
		idsData[idx] = value.ID
	}

	data := Censor{
		Content: commentsData,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		api.l.Debug().Err(err).Msg("Failed to marshal content")
		return nil, nil, err
	}

	return jsonData, idsData, nil
}

func (api *API) aiRequest(data []byte) ([]byte, error) {
	reqURL := deepPavlovModelURL

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(data))
	if err != nil {
		api.l.Debug().Err(err).Str("address", reqURL).Msg("Failed to create request")
		return nil, err
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		api.l.Debug().Err(err).Str("address", reqURL).Msg("Failed to perform POST request")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		api.l.Debug().Str("url", reqURL).Msgf("Unexpected status code: %v", response.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %v", response.StatusCode)
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		api.l.Debug().Err(err).Str("url", reqURL).Msg("Failed to read response body")
		return nil, err
	}

	return result, nil
}

func (api *API) parseURL(result chan<- CensorResult, errs chan<- error) {
	ticker := time.NewTicker(time.Second * time.Duration(period))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data, idsData, err := api.getDataFromStorage()
			if err != nil {
				errs <- err
				continue
			}

			res, err := api.aiRequest(data)
			if err != nil {
				errs <- err
				continue
			}

			select {
			case result <- CensorResult{Data: res, IDs: idsData}:
			case <-time.After(time.Second):
				api.l.Err(err).Msg("Timeout sending result")
			}
		}
	}
}
