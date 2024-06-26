package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"sephiroth-go/core"
	"sephiroth-go/core/auth"
	"sephiroth-go/core/log"
	"sephiroth-go/model/resp"
	"sephiroth-go/model/sys"
	"sephiroth-go/service"
	"sephiroth-go/util"
	"strconv"
	"time"
)

var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := auth.GetToken(c)
		if token == "" {
			resp.NoAuth("未登录或非法访问", c)
			c.Abort()
			return
		}
		if jwtService.IsBlacklist(token) {
			resp.NoAuth("您的帐户异地登陆或令牌失效", c)
			auth.ClearToken(c)
			c.Abort()
			return
		}
		j := auth.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, auth.TokenExpired) {
				resp.NoAuth("授权已过期", c)
				auth.ClearToken(c)
				c.Abort()
				return
			}
			resp.NoAuth(err.Error(), c)
			auth.ClearToken(c)
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}
		c.Set("claims", claims)
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			dr, _ := util.ParseDuration(core.Config.Jwt.ExpiresTime)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			auth.SetToken(c, newToken, int(dr.Seconds()))
			if core.Config.System.UseMultipoint {
				RedisJwtToken, err := jwtService.GetRedisJWT(newClaims.Username)
				if err != nil {
					log.Log.Error("get redis jwt failed", zap.Error(err))
				} else { // 当之前的取成功时才进行拉黑操作
					_ = jwtService.JsonInBlacklist(sys.JwtBlacklist{Jwt: RedisJwtToken})
				}
				// 无论如何都要记录当前的活跃状态
				_ = jwtService.SetRedisJWT(newToken, newClaims.Username)
			}
		}
		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}
	}
}
