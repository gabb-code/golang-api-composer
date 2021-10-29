package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gabb-code/golang-api-composer/pkg/dto"
	"github.com/gabb-code/golang-api-composer/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ComposerClientImpl struct{}

func NewComposerClientImpl() ComposerClient {
	return &ComposerClientImpl{}
}

func (client *ComposerClientImpl) Test(c *gin.Context, endpointsToCopmpose []dto.EndpointDto, method string, headers []string) (*dto.ComposedResponseDto, error) {

	response := &dto.ComposedResponseDto{}
	var endpointResponses []interface{}
	dataCh := make(chan map[string]interface{}, len(endpointResponses))
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(c)
	defer cancel()
	wg := new(sync.WaitGroup)
	cancellations := make([]context.CancelFunc, 0)
	errResult := new(error)

	go func() {
		for err := range errCh {
			errResult = &err
			fmt.Println(len(cancellations))
			wg.Done()
			for _, cancel := range cancellations {
				fmt.Println("making cancellation")
				cancel() // releases resources
			}

		}

	}()

	go func() {
		for endpointResponse := range dataCh {
			wgLock := sync.Mutex{}
			wgLock.Lock()
			endpointResponses = append(endpointResponses, endpointResponse)
			wg.Done()
			wgLock.Unlock()
			log.Printf("[INFO] [%+v]", len(endpointResponses))
		}
		fmt.Println("SENDING NIL ERROR")
	}()

	for _, endpoint := range endpointsToCopmpose {
		cancellations = append(cancellations, cancel)
		wg.Add(1)
		go DoRequest(ctx, wg, cancel, response, endpoint, dataCh, errCh)
	}
	log.Printf("Waiting...\n")
	wg.Wait()
	log.Printf("Done...\n")
	response.ComposedData = endpointResponses

	return response, *errResult
}

func DoRequest(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, response *dto.ComposedResponseDto, endpoint dto.EndpointDto, dataCh chan map[string]interface{}, errCh chan error) {
	var endpointResponse map[string]interface{}

	var payload []byte
	var err error
	if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "PATCH" {
		payload, err = json.Marshal(endpoint.Payload)
		if err != nil {
			errCh <- err
			return
		}
	}
	resp, err := utils.HTTPRequest(ctx, endpoint.URL, endpoint.Method, payload, endpoint.ContentType, endpoint.Headers)
	if err != nil {
		log.Printf("[ERR] [%v]", err)
		errCh <- err
		return
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("handler returned wrong status code: got %v instead of %v in endpoint: %v", resp.StatusCode, http.StatusOK, endpoint.URL)
		log.Printf("[ERR] [%v]", err)
		errCh <- err
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&endpointResponse); err != nil {
		log.Printf("[ERR] [%v]", err)
		errCh <- err
		return
	}

	defer resp.Body.Close()
	select {
	case <-ctx.Done():
		fmt.Printf("CONTEXT CANCELLED: %v\n", endpoint.URL)
		wg.Done()
		return
	default:
		endpointResponse["from_url"] = endpoint.URL
		dataCh <- endpointResponse
		// wg.Done()
		return
	}
}
