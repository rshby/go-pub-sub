package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-pub-sub/enternal/httpclient"
	"io"
	"net/http"
	"testing"
)

func TestDecodeRequestBody(t *testing.T) {
	t.Run("decode request body nil", func(t *testing.T) {
		req := map[string]any{}

		requestBody, err := json.Marshal(&req)
		assert.Nil(t, err)
		fmt.Println(string(requestBody))
	})
}

func TestRequestHTTP(t *testing.T) {
	client := httpclient.NewHttpClient()

	// create request
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:4002/v1/memberships/customer-vertical/12336", nil)
	assert.Nil(t, err)

	req.Header.Add("Authorization", "Basic ZW1lbWJlcnBvaW50YmFzaWM6Y2hQWEJac2Jueg==")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBodyBytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := map[string]any{}
	err = json.Unmarshal(resBodyBytes, &responseBody)
	assert.Nil(t, err)

	responseBodyJson, err := json.Marshal(&responseBody)
	assert.Nil(t, err)

	fmt.Println(string(responseBodyJson))
}
