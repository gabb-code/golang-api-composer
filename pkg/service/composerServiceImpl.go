package service

import (
	"log"
	"net/http"

	"github.com/gabb-code/golang-api-composer/pkg/client"
	"github.com/gabb-code/golang-api-composer/pkg/dto"
	"github.com/gin-gonic/gin"
)

type ComposerServiceImpl struct {
	client client.ComposerClient
}

func NewComposerServiceImpl(client client.ComposerClient) ComposerService {
	return &ComposerServiceImpl{
		client: client,
	}
}

func (service ComposerServiceImpl) Test(c *gin.Context) {

	var endpointsToCopmpose []dto.EndpointDto
	var resp *dto.ComposedResponseDto
	var err error

	if err := c.BindJSON(&endpointsToCopmpose); err != nil {
		log.Printf("[ERR] [%v]", err)
		resp.StatusCode = 400
		resp.StatusDesc = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if resp, err = service.client.Test(c, endpointsToCopmpose, "GET", nil); err != nil {
		log.Printf("[ERR] [%v]", err)
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusDesc = err.Error()
		resp.ComposedData = nil
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.StatusCode = http.StatusOK
	resp.StatusDesc = "Success"
	resp.Success = true
	c.JSON(200, resp)

}
