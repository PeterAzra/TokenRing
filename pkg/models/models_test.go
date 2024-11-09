package models

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnmarshalJoinResponse(t *testing.T) {
	jsonData := []byte(`{"Ok":true,"Left":"http://localhost:8001","Right":"http://localhost:8002"}`)
	var response JoinResponse
	valid := json.Valid(jsonData)
	assert.True(t, valid)
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	assert.True(t, response.Ok)
	assert.True(t, len(response.Left) > 0)
	assert.True(t, len(response.Right) > 0)
}
