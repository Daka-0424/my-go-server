package route

import (
	"net/http"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/web/admin"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/gin-gonic/gin"
)

func AdminRoute(
	route *gin.Engine,
	cfg *config.Config,
	admin *admin.AdminController,
) {
	logoutCheckGroup := route.Group("/admin", checkLogout(admin))
	if cfg.IsDevelopment() {
		logoutCheckGroup.POST("/register", admin.PostTempRegisterRequest)
	}
}

func checkLogin(ac *admin.AdminController, roleType entity.AdminRoleType) gin.HandlerFunc {
	return func(c *gin.Context) {
		admin := ac.GetSession(c)
		if admin == nil {
			c.Redirect(http.StatusFound, "/admin/signup")
			c.Abort()
		} else {
			if checkAdminRole(admin, roleType) {
				c.Next()
			} else {
				c.Redirect(http.StatusFound, "/admin")
				c.Abort()
			}
		}
	}
}

func checkAdminRole(admin *entity.Admin, role entity.AdminRoleType) bool {
	switch admin.RoleType {
	case entity.AdminRoleTypeMaster:
		return admin.RoleType == role
	case entity.AdminRoleTypeGuest:
		return admin.RoleType != role
	default:
		return true
	}
}

func checkLogout(ac *admin.AdminController) gin.HandlerFunc {
	return func(c *gin.Context) {
		admin := ac.GetSession(c)
		if admin != nil {
			c.Redirect(http.StatusFound, "/admin")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
