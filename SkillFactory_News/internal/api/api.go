package api

import (
	"context"
	"net/http"
	"time"

	storage "github.com/MaksimovDenis/SkillFactory_News/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 //512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Opts struct {
	Addr    string
	Log     zerolog.Logger
	Storage *storage.Storage
}

type API struct {
	l       zerolog.Logger
	server  *http.Server
	router  *gin.Engine
	storage *storage.Storage
}

func NewAPI(opts *Opts) *API {
	router := gin.Default()

	api := &API{
		l: opts.Log,
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		router:  router,
		storage: opts.Storage,
	}

	go api.StartParseUrl()

	api.setupEndpoints()

	return api
}

func (api *API) setupEndpoints() {
	api.router.GET("api/feeds/:limit", api.Feeds)
	api.router.GET("api/feed/:id", api.FeedsById)
	api.router.POST("api/feeds", api.FeedsByFilter)
}

func (api *API) Serve() error {
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to start api server")
		return err
	}
	return nil
}

func (api *API) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, ServerShutdownDefaultDelay)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to stop api server")
	}
}
