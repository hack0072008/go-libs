package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Middleware.CORS
 */
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers",
			"X-UCLOUD-REQUEST-ID, Content-Type, Content-Length, Accept-Encoding,Authorization, accept, origin, Cache-Control, X-Requested-With, From")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		ctx.Next()
	}
}
