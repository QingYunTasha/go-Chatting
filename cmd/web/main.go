package main

import (
	config "go-Chatting/config"
	webdelivery "go-Chatting/internal/app/web/delivery"
	webinfra "go-Chatting/internal/app/web/infra"
	webusecase "go-Chatting/internal/app/web/usecase"
	dbfactory "go-Chatting/internal/infra/database/factory"
	utils "go-Chatting/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.Init()

	db, err := dbfactory.InitDb(viper.GetString("DB_DSN"))
	if err != nil {
		panic(err.Error())
	}

	dbRepository, err := dbfactory.NewDbRepository(db)
	if err != nil {
		panic(err.Error())
	}

	mailer := webinfra.NewSMTPMailer(
		viper.GetString("mailer.host"),
		viper.GetInt("mailer.port"),
		viper.GetString("mailer.username"),
		viper.GetString("mailer.password"),
	)

	secretKey := utils.NewSecretKey()

	webUsecase := webusecase.NewWebUsecase(dbRepository, mailer, secretKey)

	server := gin.Default()

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "ok",
		})
	})

	webdelivery.NewWebHandler(server, webUsecase, secretKey)
	server.Run(":8080")
}
