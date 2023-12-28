package route

import (
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Route(
	route *gin.Engine,
	cfg *config.Config,
	cache repository.Cache,
	localizer *i18n.Localizer,
	registration *controller.RegistrationController,
) {
	route.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "ok"}) })

	// 認証なし
	route.POST("/api/registration", registration.CreateRegistration)
}
