package api

import (
	"context"
	"net/http"
	"skillfactory/finalProject/commentsService/internal/api/oapi"
	"skillfactory/finalProject/commentsService/internal/storage"
	"time"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 //512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Opts struct {
	Addr        string
	Log         zerolog.Logger
	Storage     *storage.Storage
	SwaggerFile []byte
}

type API struct {
	l           zerolog.Logger
	server      *http.Server
	router      *gin.Engine
	storage     *storage.Storage
	swaggerFile []byte
}

func NewAPI(opts *Opts) (*API, error) {
	router := gin.Default()

	swagger, err := oapi.GetSwagger()
	if err != nil {
		return nil, err
	}

	oapiOpts := &middleware.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody: true,
		},
	}

	router = gin.Default()

	router.MaxMultipartMemory = FileUploadBufferSize

	api := &API{
		l: opts.Log,
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		router:      router,
		storage:     opts.Storage,
		swaggerFile: opts.SwaggerFile,
	}

	router.Use(middleware.OapiRequestValidatorWithOptions(swagger, oapiOpts))

	oapi.RegisterHandlersWithOptions(router, api, oapi.GinServerOptions{
		BaseURL: "/api",
	})

	go api.StartCensorship()

	return api, nil
}

func (hdl *API) Serve() {
	go func() {
		if err := hdl.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hdl.l.Error().Err(err).Msg("failed to start api server")
		}
	}()
}

func (hdl *API) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownDefaultDelay)
	defer cancel()

	if err := hdl.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		hdl.l.Error().Err(err).Msg("failed to stop api server")
	}
}
