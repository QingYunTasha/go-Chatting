package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	server := gin.Default()

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": "ok",
		})
	})

	server.Run(":8080")
}
