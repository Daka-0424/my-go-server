package route

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/web/cheat"
	"github.com/gin-gonic/gin"
)

func CheatRoute(
	route *gin.Engine,
	cfg *config.Config,
	cheatRoute *cheat.CheatRootController,
) {
	route.LoadHTMLGlob("pkg/controller/web/**/**/**/*.html")
	route.Static("/assets", "assets")

	// 開発環境以外では表示しない
	if !cfg.IsDevelopment() {
		return
	}

	cheatRoot := route.Group("/cheat")
	cheatRoot.GET("/", cheatRoute.Get)
}
