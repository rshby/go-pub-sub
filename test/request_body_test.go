package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type RequestBody struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
}

func TestGetRequestBody(t *testing.T) {
	req := map[string]any{
		"phone_number": "08122222",
	}
	reqJson, _ := json.Marshal(&req)

	// decode to request
	var request RequestBody
	json.Unmarshal(reqJson, &request)

	fmt.Println(request)
}

func TestRequestBodyWithNull(t *testing.T) {
	requestBody := `{
		"phone_number" : "083863890419",
		"email" : null 
	}`

	var req RequestBody
	err := json.Unmarshal([]byte(requestBody), &req)
	assert.Nil(t, err)

	fmt.Println(req)

	reqJson, _ := json.Marshal(&req)
	fmt.Println(string(reqJson))

}
