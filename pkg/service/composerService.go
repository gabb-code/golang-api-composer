package service

import "github.com/gin-gonic/gin"

type ComposerService interface {
	Test(c *gin.Context)
}
