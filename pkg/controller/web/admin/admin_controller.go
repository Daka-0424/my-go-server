package admin

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AdminController struct {
	adminControllerBase
	adminRepository repository.Admin
	adminUsecase    usecase.Admin
}

func NewAdminController(
	cfg *config.Config,
	lc *i18n.Localizer,
	cache repository.Cache,
	ar repository.Admin,
	au usecase.Admin,
) *AdminController {
	return &AdminController{
		adminControllerBase: adminControllerBase{
			cfg:       cfg,
			localizer: lc,
			cache:     cache,
		},
		adminRepository: ar,
		adminUsecase:    au,
	}
}

func (ctl *AdminController) PostRegisterRequest(ctx *gin.Context) {
	email := ctx.PostForm("email")

	if ctl.adminRepository.Exsists(ctx, email) {
		ctx.Redirect(302, "/admin/register")
		return
	}

	key := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic("ランダムな文字列の生成に失敗しました。")
	}
	redisKey := base64.URLEncoding.EncodeToString(key)

	ctl.cache.Set(ctx, redisKey, []byte(email), time.Hour)

	// 一旦リダイレクト
	ctx.Redirect(302, "/admin/register/"+redisKey)

	// TODO:メール送信処理
}

func (ctl *AdminController) PostRegister(ctx *gin.Context) {
	key := ctx.Param("value")

	email, ok, err := ctl.cache.Get(ctx, key)
	// 管理者登録の有効期間切れ等
	if err != nil || !ok {
		ctx.Redirect(302, "/admin/register")
		return
	}

	pass := ctx.PostForm("password")

	admin, err := ctl.adminUsecase.Register(ctx, string(email), pass, entity.AdminRoleTypeBasic)
	// 管理者登録失敗
	if err != nil {
		ctx.Redirect(302, "/admin/register")
		return
	}

	ctl.newSession(ctx, admin)
	ctx.Redirect(301, "/admin")
}

func (ctl *AdminController) newSession(ctx *gin.Context, admin *entity.Admin) {
	baseRedisKey := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, baseRedisKey); err != nil {
		panic("ランダムな文字列の生成に失敗しました。")
	}
	newRedisKey := base64.URLEncoding.EncodeToString(baseRedisKey)

	data := *admin
	data.Password = ""
	json, err := json.Marshal(data)
	if err != nil {
		panic("セッション情報の生成に失敗しました。")
	}

	key := ctl.cfg.Cookie.Key
	host := ctl.cfg.Cookie.Host
	ctl.cache.Set(ctx, newRedisKey, json, time.Hour*10)
	ctx.SetCookie(key, newRedisKey, 3600*10, "/", host, false, true)
}
