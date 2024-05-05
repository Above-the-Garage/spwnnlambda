package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	InitDict()
	var pathParam = make(map[string]string)
	result, err := Handler(events.APIGatewayProxyRequest{PathParameters: pathParam})
	assert.Equal(t, ErrNoWordToCheck, err)
	pathParam["word"] = "stephen"
	result, err = Handler(events.APIGatewayProxyRequest{PathParameters: pathParam})
	assert.IsType(t, nil, err)
	assert.NotEqual(t, 0, len(result.Body))
}
