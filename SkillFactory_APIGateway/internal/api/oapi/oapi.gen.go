// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package oapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// Comment defines model for Comment.
type Comment struct {
	// Content Content of the comment
	Content string `json:"content"`

	// Id ID of the comment
	Id int `json:"id"`

	// NewsId ID of the associated news item
	NewsId int `json:"news_id"`

	// ParentCommentId ID of the parent comment (if any)
	ParentCommentId *int `json:"parent_comment_id,omitempty"`
}

// CreateComment defines model for CreateComment.
type CreateComment struct {
	// Content Content of the comment
	Content string `json:"content"`

	// NewsId ID of the associated news item
	NewsId int `json:"news_id"`

	// ParentCommentId ID of the parent comment (if any)
	ParentCommentId *int `json:"parent_comment_id"`
}

// Feeds defines model for Feeds.
type Feeds struct {
	// Content Content of the feeds
	Content string `json:"content"`

	// Id Feeds ID
	Id int `json:"id"`

	// Link Link of the feeds
	Link string `json:"link"`

	// PubDate Publication date
	PubDate string `json:"pub_date"`

	// Title Title of feeds (if any)
	Title string `json:"title"`
}

// FeedsByFilter defines model for FeedsByFilter.
type FeedsByFilter struct {
	// Filter Filter contains string that will find all comments which title or content include this string
	Filter string `json:"filter"`

	// Limit Length of feeds list
	Limit int `json:"limit"`
}

// ID defines model for ID.
type ID = int64

// Limit defines model for Limit.
type Limit = int64

// FeedsParams defines parameters for Feeds.
type FeedsParams struct {
	// Limit maximum number of results to return
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`
}

// CreateCommentJSONRequestBody defines body for CreateComment for application/json ContentType.
type CreateCommentJSONRequestBody = CreateComment

// FeedsByFilterJSONRequestBody defines body for FeedsByFilter for application/json ContentType.
type FeedsByFilterJSONRequestBody = FeedsByFilter

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetAllComments request
	GetAllComments(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateCommentWithBody request with any body
	CreateCommentWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateComment(ctx context.Context, body CreateCommentJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Feeds request
	Feeds(ctx context.Context, params *FeedsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// FeedsByFilterWithBody request with any body
	FeedsByFilterWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	FeedsByFilter(ctx context.Context, body FeedsByFilterJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// FeedsById request
	FeedsById(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSwagger request
	GetSwagger(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAllComments(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAllCommentsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateCommentWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateCommentRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateComment(ctx context.Context, body CreateCommentJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateCommentRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Feeds(ctx context.Context, params *FeedsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewFeedsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) FeedsByFilterWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewFeedsByFilterRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) FeedsByFilter(ctx context.Context, body FeedsByFilterJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewFeedsByFilterRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) FeedsById(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewFeedsByIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetSwagger(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetSwaggerRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAllCommentsRequest generates requests for GetAllComments
func NewGetAllCommentsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/comments")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateCommentRequest calls the generic CreateComment builder with application/json body
func NewCreateCommentRequest(server string, body CreateCommentJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateCommentRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateCommentRequestWithBody generates requests for CreateComment with any type of body
func NewCreateCommentRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/comments")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewFeedsRequest generates requests for Feeds
func NewFeedsRequest(server string, params *FeedsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/feeds")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewFeedsByFilterRequest calls the generic FeedsByFilter builder with application/json body
func NewFeedsByFilterRequest(server string, body FeedsByFilterJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewFeedsByFilterRequestWithBody(server, "application/json", bodyReader)
}

// NewFeedsByFilterRequestWithBody generates requests for FeedsByFilter with any type of body
func NewFeedsByFilterRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/feeds")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewFeedsByIdRequest generates requests for FeedsById
func NewFeedsByIdRequest(server string, id ID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/feeds/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetSwaggerRequest generates requests for GetSwagger
func NewGetSwaggerRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/swagger")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetAllCommentsWithResponse request
	GetAllCommentsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllCommentsResponse, error)

	// CreateCommentWithBodyWithResponse request with any body
	CreateCommentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateCommentResponse, error)

	CreateCommentWithResponse(ctx context.Context, body CreateCommentJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateCommentResponse, error)

	// FeedsWithResponse request
	FeedsWithResponse(ctx context.Context, params *FeedsParams, reqEditors ...RequestEditorFn) (*FeedsResponse, error)

	// FeedsByFilterWithBodyWithResponse request with any body
	FeedsByFilterWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*FeedsByFilterResponse, error)

	FeedsByFilterWithResponse(ctx context.Context, body FeedsByFilterJSONRequestBody, reqEditors ...RequestEditorFn) (*FeedsByFilterResponse, error)

	// FeedsByIdWithResponse request
	FeedsByIdWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*FeedsByIdResponse, error)

	// GetSwaggerWithResponse request
	GetSwaggerWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSwaggerResponse, error)
}

type GetAllCommentsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Comment
}

// Status returns HTTPResponse.Status
func (r GetAllCommentsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAllCommentsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateCommentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r CreateCommentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateCommentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type FeedsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Feeds
}

// Status returns HTTPResponse.Status
func (r FeedsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FeedsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type FeedsByFilterResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Feeds
}

// Status returns HTTPResponse.Status
func (r FeedsByFilterResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FeedsByFilterResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type FeedsByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Feeds
}

// Status returns HTTPResponse.Status
func (r FeedsByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FeedsByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSwaggerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetSwaggerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSwaggerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAllCommentsWithResponse request returning *GetAllCommentsResponse
func (c *ClientWithResponses) GetAllCommentsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllCommentsResponse, error) {
	rsp, err := c.GetAllComments(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAllCommentsResponse(rsp)
}

// CreateCommentWithBodyWithResponse request with arbitrary body returning *CreateCommentResponse
func (c *ClientWithResponses) CreateCommentWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateCommentResponse, error) {
	rsp, err := c.CreateCommentWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateCommentResponse(rsp)
}

func (c *ClientWithResponses) CreateCommentWithResponse(ctx context.Context, body CreateCommentJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateCommentResponse, error) {
	rsp, err := c.CreateComment(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateCommentResponse(rsp)
}

// FeedsWithResponse request returning *FeedsResponse
func (c *ClientWithResponses) FeedsWithResponse(ctx context.Context, params *FeedsParams, reqEditors ...RequestEditorFn) (*FeedsResponse, error) {
	rsp, err := c.Feeds(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFeedsResponse(rsp)
}

// FeedsByFilterWithBodyWithResponse request with arbitrary body returning *FeedsByFilterResponse
func (c *ClientWithResponses) FeedsByFilterWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*FeedsByFilterResponse, error) {
	rsp, err := c.FeedsByFilterWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFeedsByFilterResponse(rsp)
}

func (c *ClientWithResponses) FeedsByFilterWithResponse(ctx context.Context, body FeedsByFilterJSONRequestBody, reqEditors ...RequestEditorFn) (*FeedsByFilterResponse, error) {
	rsp, err := c.FeedsByFilter(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFeedsByFilterResponse(rsp)
}

// FeedsByIdWithResponse request returning *FeedsByIdResponse
func (c *ClientWithResponses) FeedsByIdWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*FeedsByIdResponse, error) {
	rsp, err := c.FeedsById(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFeedsByIdResponse(rsp)
}

// GetSwaggerWithResponse request returning *GetSwaggerResponse
func (c *ClientWithResponses) GetSwaggerWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSwaggerResponse, error) {
	rsp, err := c.GetSwagger(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSwaggerResponse(rsp)
}

// ParseGetAllCommentsResponse parses an HTTP response from a GetAllCommentsWithResponse call
func ParseGetAllCommentsResponse(rsp *http.Response) (*GetAllCommentsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAllCommentsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Comment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateCommentResponse parses an HTTP response from a CreateCommentWithResponse call
func ParseCreateCommentResponse(rsp *http.Response) (*CreateCommentResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateCommentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseFeedsResponse parses an HTTP response from a FeedsWithResponse call
func ParseFeedsResponse(rsp *http.Response) (*FeedsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &FeedsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Feeds
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseFeedsByFilterResponse parses an HTTP response from a FeedsByFilterWithResponse call
func ParseFeedsByFilterResponse(rsp *http.Response) (*FeedsByFilterResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &FeedsByFilterResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Feeds
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseFeedsByIdResponse parses an HTTP response from a FeedsByIdWithResponse call
func ParseFeedsByIdResponse(rsp *http.Response) (*FeedsByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &FeedsByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Feeds
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetSwaggerResponse parses an HTTP response from a GetSwaggerWithResponse call
func ParseGetSwaggerResponse(rsp *http.Response) (*GetSwaggerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetSwaggerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /comments)
	GetAllComments(c *gin.Context)

	// (POST /comments)
	CreateComment(c *gin.Context)

	// (GET /feeds)
	Feeds(c *gin.Context, params FeedsParams)

	// (POST /feeds)
	FeedsByFilter(c *gin.Context)

	// (GET /feeds/{id})
	FeedsById(c *gin.Context, id ID)

	// (GET /swagger)
	GetSwagger(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetAllComments operation middleware
func (siw *ServerInterfaceWrapper) GetAllComments(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetAllComments(c)
}

// CreateComment operation middleware
func (siw *ServerInterfaceWrapper) CreateComment(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateComment(c)
}

// Feeds operation middleware
func (siw *ServerInterfaceWrapper) Feeds(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params FeedsParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter limit: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Feeds(c, params)
}

// FeedsByFilter operation middleware
func (siw *ServerInterfaceWrapper) FeedsByFilter(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.FeedsByFilter(c)
}

// FeedsById operation middleware
func (siw *ServerInterfaceWrapper) FeedsById(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id ID

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.FeedsById(c, id)
}

// GetSwagger operation middleware
func (siw *ServerInterfaceWrapper) GetSwagger(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetSwagger(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/comments", wrapper.GetAllComments)
	router.POST(options.BaseURL+"/comments", wrapper.CreateComment)
	router.GET(options.BaseURL+"/feeds", wrapper.Feeds)
	router.POST(options.BaseURL+"/feeds", wrapper.FeedsByFilter)
	router.GET(options.BaseURL+"/feeds/:id", wrapper.FeedsById)
	router.GET(options.BaseURL+"/swagger", wrapper.GetSwagger)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xXTY/bNhD9KwTbQwuoltO0Pei2WWMDA1u0QHsLFgYtjayJKVIhR7s1DP33gqQoS7G8",
	"6xRIsxcDIufjzZunGfnIc103WoEiy7Mjb4QRNRAY/7Reud8CbG6wIdSKZxwLpkumtx8hJ55wdGeNoIon",
	"XIkavAVPuIFPLRooeEamhYTbvIJauHClNrUgZ6fot194wunQQHiEHRjedQm/xxrpPHct/sG6rZlq6y0Y",
	"h8OAbSVZRpoZoNaoCOlTC+ZwwiR9wDGMAkrRSuLZm2VyBaYu+npibnVdg/IIG6MbMITgL3KtqL+YQr8N",
	"Fw4zVcDyPsCQyZJBteNd4ug7816vLjoOEBOu4MlunncX1uocBUHBnDVDgno2ViMMKNr06V6IGowjNvYD",
	"lkyow4/zvT0p40OQSoSdDOw9DI69zrqE3xoQBF+T+FdIn2qlFFsJ8S16gc7rmLwDKOx/Z7D07lcK1+di",
	"69UsSRLV/tzlHtX+xVxNu90UguDc/c92KzEX7ol5ixlvQpIzrn+7Y5fap52RcYwwp+IQ88R8X94I6sVe",
	"vDvcoSQw5z0ph/PPaPXnzOUSqCwLwBhVgtgTSslKVAUTUkZRWfZUYV4xCiUGV9dUVLlsC2BUYQwzx5ic",
	"H8j3oHZUnSiTaOnltz4O4766c16cA6pSR1mK3OdujeQZr4gam6Wp3aOUpchJm8PCtKmDOUX3/o+f3v7K",
	"fhd7i7V+ZCtQaIdOZXwcYFOiEnLTGN3vtUcwNkRZLpaLNy64bkCJBnnG3y6Wi6VrraDK9ymNNLuHHcww",
	"9R5o0o8F9wGNF+q6CBY3Ut7GQI4y22hlgxJ+Xi4/e0lF00Shpx+tS3IcrTc3mbzj9wZKnvHv0tOiT/tN",
	"lsZp2g0dEMaIQ2jAFP+N763r9FCqMyKxs66lw+GDezm1nRsjfoBH9/P6pwM+KAYsvdPF4Ysqf7bgSY5u",
	"Kkw3Yrt52qeV6P2l2ruEp2Ucr8/qwFudk3DXD7zxd9iH+ZpOJmn4XOoe/g/NBIRXKOa+10tf6ZixQNFl",
	"qYStsT2wMCEu0DSMza+jlWmO67Xyjei+8WQZKPqlOcf3IM/0iEX3rEbL2AIsLtK/Lr5YqevVa5Pp3bjQ",
	"S6TZJ7HbhUU8y9hf4d61AOYGe3/Pr5ku+elzy47DjrE1Rm/BYXN/S8A8RvbDhkzdlnI09+bH+C+oV0Yy",
	"HJxm+UP3bwAAAP//uCW3kgsOAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
