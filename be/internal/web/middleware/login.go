package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	IgnorePaths map[string]bool
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		IgnorePaths: make(map[string]bool),
	}
}

func (b *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	b.IgnorePaths[path] = true
	return b
}

func (b *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.FullPath()

		// 不需要登录验证
		if _, ok := b.IgnorePaths[path]; ok {
			return
		}

		// 获取 session 里有无userId,
		// 有： 登录过， 没有： 没登录
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新token 过期时间
		updateTime := sess.Get("update_time")
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		if updateTime == nil {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

	}
}
