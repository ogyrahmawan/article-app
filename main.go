package main

import (
	"article-app/internal"
	"article-app/internal/middlewares"
	"strings"

	userHandler "article-app/internal/data/user/delivery/http"
	userRepo "article-app/internal/data/user/repository"
	userUsecase "article-app/internal/data/user/usecase"

	articleHandler "article-app/internal/data/article/delivery/http"
	articleRepo "article-app/internal/data/article/repository"
	articleUsecase "article-app/internal/data/article/usecase"
	"article-app/internal/domain"
	"article-app/pkg/database"
	"article-app/pkg/jwt"
	"article-app/pkg/seeder"
	"errors"
	"log"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/beego/i18n"
	"gorm.io/gorm"
)

func main() {

	// token expired
	tokenExpired := beego.AppConfig.DefaultInt64("tokenExpired", 86400)
	// global execution timeout
	serverTimeout := beego.AppConfig.DefaultInt64("serverTimeout", 60)
	// global execution timeout
	requestTimeout := beego.AppConfig.DefaultInt("executionTimeout", 5)
	// global execution timeout to second
	timeoutContext := time.Duration(requestTimeout) * time.Second
	// jwt secret key
	jwtSecretKey := beego.AppConfig.DefaultString("jwtSecretKey", "secret")
	// log path

	// languange
	lang := beego.AppConfig.DefaultString("lang", "en|id")
	languages := strings.Split(lang, "|")
	for _, value := range languages {
		if err := i18n.SetMessage(value, "conf/"+value+".ini"); err != nil {
			panic("Failed to set message file for l10n")
		}
	}

	// beego config
	beego.BConfig.Log.AccessLogs = false
	beego.BConfig.Log.EnableStaticLogs = false
	beego.BConfig.Listen.ServerTimeOut = serverTimeout

	// database initialization
	db := database.DB()

	if beego.BConfig.RunMode != "prod" {
		// db auto migrate dev environment
		err := db.AutoMigrate(
			&domain.User{},
			&domain.Article{},
		)
		if err == nil && db.Migrator().HasTable(&domain.User{}) {
			if err := db.First(&domain.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				//Insert seed data
				seeder.Seeds(db)
			}
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			panic(err)
		}
	}

	beego.BeeApp.Server.RegisterOnShutdown(func() {
		if sqlDb, err := db.DB(); err != nil {
			log.Println("error database connection ...")
		} else {
			sqlDb.Close()
			log.Println("close database connection ...")
		}
	})

	// health check
	beego.Get("/health", func(ctx *beegoContext.Context) {
		ctx.Output.SetStatus(http.StatusOK)
		ctx.Output.JSON(beego.M{"status": "alive"}, beego.BConfig.RunMode != "prod", false)
	})

	// jwt middleware
	auth, err := jwt.NewJwt(&jwt.Options{
		SignMethod:  jwt.HS256,
		SecretKey:   jwtSecretKey,
		Locations:   "header:Authorization",
		IdentityKey: "uid",
	})
	if err != nil {
		panic(err)
	}

	// middleware init
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
	}))
	beego.InsertFilterChain("*", middlewares.RequestID())
	beego.InsertFilterChain("/api/v1/*", middlewares.NewJwtMiddleware().JwtMiddleware(auth))

	// default error handler
	beego.ErrorController(&internal.BaseController{})

	// init repository
	userRepository := userRepo.NewUserRepository(db)
	articleRepository := articleRepo.NewArticleRepository(db)

	// init usecase
	userUsecase := userUsecase.NewUserUseCase(timeoutContext, userRepository, auth, int(tokenExpired))
	articleUsecase := articleUsecase.NewArticleUseCase(timeoutContext, articleRepository, auth, int(tokenExpired))

	// init handler
	userHandler.NewUserHandler(userUsecase, auth)
	articleHandler.NewArticleHandler(articleUsecase, auth)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	beego.Run()

}
