package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
	"webook/be/config"
	"webook/be/internal/repository"
	"webook/be/internal/repository/dao"
	"webook/be/internal/service"
	"webook/be/internal/web"
	"webook/be/internal/web/middleware"
	"webook/be/pkg/ginx/middlewares/ratelimit"
)

func main() {
	db := initDB()
	redisClient := initRedis()

	server := initServer(redisClient)

	initUser(db, server)

	if err := server.Run(":8080"); err != nil {
		log.Println(err)
	}
}

func initRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{Addr: config.Config.Redis.Addr})
}

func initUser(db *gorm.DB, server *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userSvc := service.NewUserService(userRepo)
	userHdl := web.NewUserHandler(userSvc)
	userHdl.RegisterRoutes(server)
}

func initServer(redisClient redis.Cmdable) *gin.Engine {
	server := gin.Default()

	// 限流
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	// 跨域
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,                                                                  // 是否允许带 cookie
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // 允许的请求头
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool { // 哪些来源的url是被允许的
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour, // preflight响应 过期时间
	}))

	server.Use(middleware.NewLoginJWTMiddlewareBuilder(web.UserJWTSignedString).IgnorePath("/users/signup").IgnorePath("/users/login").Build())

	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(err)
	}
	return db
}
