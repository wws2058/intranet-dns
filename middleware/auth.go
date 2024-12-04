package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/intranet-dns/ctx"
	"github.com/tswcbyy1107/intranet-dns/models"
	"github.com/tswcbyy1107/intranet-dns/service"
	"github.com/tswcbyy1107/intranet-dns/utils"
)

var whiteApis = []string{
	"/api/v1/users/login",
	"/api/v1/ping",
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip auth
		if utils.Contains[string](whiteApis, c.FullPath()) || strings.Contains(c.FullPath(), "swagger") {
			c.Next()
			return
		}

		// check if jwt token legal
		jwtToken := c.Request.Header.Get("token")
		claim, err := service.ParseToken(jwtToken)
		if err != nil {
			ctx.AbortRsp(c, err)
			return
		}

		// check user api permission
		user := &models.SysUser{
			Name: claim.Username,
		}
		err = models.TemplateQuery(&models.DaoDBReq{
			Dst:         user,
			ModelFilter: user,
		})
		if err != nil {
			ctx.AbortRsp(c, err)
			return
		}
		ctx.SetLoginUsername(c, user.Name)
		if !user.Active {
			ctx.AbortRsp(c, fmt.Errorf("user is inactive"))
			return
		}
		// super admin
		if user.IsSuperAdmin() {
			c.Next()
			go models.UpdateUserLoginInfo(user.Name)
			return
		}
		// exclude inactive apis
		permissions, _ := user.GetAccessibleApis()
		if !utils.Contains[string](permissions, fmt.Sprintf("%s/%s", c.FullPath(), c.Request.Method)) {
			ctx.AbortRsp(c, fmt.Errorf("has no %s/%s permission", c.FullPath(), c.Request.Method))
			return
		}

		c.Next()
		go models.UpdateUserLoginInfo(user.Name)
	}
}
