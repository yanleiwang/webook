package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/be/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
	IgnorePaths  []string
	SignedString []byte
}

func NewLoginJWTMiddlewareBuilder(signedString []byte) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		SignedString: signedString,
	}

}

func (b *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {

	b.IgnorePaths = append(b.IgnorePaths, path)
	return b
}

func (b *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		for _, val := range b.IgnorePaths {
			if path == val {
				return
			}
		}
		// Create claims while leaving out some of the optional fields
		//claims = MyCustomClaims{
		//	"bar",
		//	jwt.RegisteredClaims{
		//		// Also fixed dates can be used for the NumericDate
		//		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
		//		Issuer:    "test",
		//	},
		//}
		//
		//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		//ss, err := token.SignedString(mySigningKey)
		//fmt.Printf("%v %v", ss, err)

		tokenHeader := ctx.GetHeader("Auth")
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录，有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return b.SignedString, nil
		})

		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//claims.ExpiresAt.Time.Before(time.Now()) {
		//	// 过期了
		//}
		// err 为 nil，token 不为 nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		// 每十秒钟刷新一次
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(b.SignedString)
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)

	}
}
