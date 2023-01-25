package internal

import "github.com/gin-gonic/gin"

func NewHttpServer() *gin.Engine {
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		if ctx.Keys == nil {
			ctx.Keys = map[string]any{}
		}
		ctx.Next()
	})
	return r
}
