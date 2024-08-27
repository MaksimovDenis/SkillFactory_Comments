package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) GetSwagger(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, api.swaggerFile)
}
