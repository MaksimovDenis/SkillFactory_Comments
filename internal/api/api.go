package api

import (
	"skillfactory/finalProject/commentsService/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Opts struct {
	Addr    string
	Log     zerolog.Logger
	Storage *storage.Storage
}

type API struct {
	l       zerolog.Logger
	router  *gin.Engine
	storage *storage.Storage
}
