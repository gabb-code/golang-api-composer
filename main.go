package main

import (
	"log"

	"github.com/gabb-code/golang-api-composer/pkg/client"
	"github.com/gabb-code/golang-api-composer/pkg/service"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {

	gin.SetMode("release")
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	composerClient := client.NewComposerClientImpl()
	composerService := service.NewComposerServiceImpl(composerClient)

	r.POST("/test", composerService.Test)
}

func main() {
	if err := r.Run(":6000"); err != nil {
		log.Fatal(err)
	}
}
