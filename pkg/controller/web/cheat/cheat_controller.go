package cheat

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/gin-gonic/gin"
)

type CheatRootController struct {
	cfg *config.Config
}

func NewCheatRootController(cfg *config.Config) *CheatRootController {
	return &CheatRootController{
		cfg: cfg,
	}
}

func (c *CheatRootController) Get(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "cheat/index", gin.H{})
}
