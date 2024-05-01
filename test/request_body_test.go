package test

import (
	"encoding/json"
	"fmt"
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
