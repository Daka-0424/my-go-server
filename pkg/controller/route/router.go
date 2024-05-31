package route

import (
	"github.com/Daka-0424/my-go-server/config"
	controller "github.com/Daka-0424/my-go-server/pkg/controller/api"
	"github.com/Daka-0424/my-go-server/pkg/controller/middleware"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Route(
	route *gin.Engine,
	cfg *config.Config,
	cache repository.ICache,
	localizer *i18n.Localizer,
	registration *controller.UserController,
	session *controller.SessionController,
	vcPlatformProduct *controller.PlatformProductController,
	appstore *controller.AppstoreController,
	playstore *controller.PlaystoreController,
) {
	route.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "ok"}) })

	// 認証なし
	route.POST("/api/registration", registration.Registration)
	route.POST("/api/session", session.CreateSession)

	authMiddleware := middleware.JwtMiddleware(cfg, localizer, cache)

	sessionGroup := route.Group("api/my").Use(authMiddleware)
	sessionGroup.GET("/platform-products", vcPlatformProduct.ListPlatformProduct)

	sessionGroup.POST("/appstore/billing", appstore.Billing)
	sessionGroup.POST("/playstore/billing", playstore.Billing)
}
