package client

import (
	"github.com/gabb-code/golang-api-composer/pkg/dto"
	"github.com/gin-gonic/gin"
)

type ComposerClient interface {
	Test(c *gin.Context, cndpointsToCopmpose []dto.EndpointDto, method string, headers []string) (*dto.ComposedResponseDto, error)
}
