package middleware

import (
	"errors"
	"strings"

	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/gin-gonic/gin"
)

func returnErrorWithAbort(ctx *gin.Context, localizer *language.Localizer) {
	lang := ctx.Request.Header.Get("Accept-Language")
	appErr := response.NewErrUnauthorized(model.E0101, localizer.MustLocalize(model.E0101, lang, nil))
	formatter.Respond(ctx, appErr.StatusCode, gin.H{"error": appErr})
	ctx.Abort()
}

func returnErrorWithAbortForMultiDevice(ctx *gin.Context, localizer *language.Localizer) {
	lang := ctx.Request.Header.Get("Accept-Language")
	appErr := response.NewErrUnauthorized(model.E9902, localizer.MustLocalize(model.E9902, lang, nil))
	formatter.Respond(ctx, appErr.StatusCode, gin.H{"error": appErr})
	ctx.Abort()
}

func bearerToken(ctx *gin.Context) (string, error) {
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("token not found")
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		return "", errors.New("token not found")
	}

	return token, nil
}
