package route

import "github.com/gin-gonic/gin"

func Route(
	route *gin.Engine,
) {
	route.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "ok"}) })
}
