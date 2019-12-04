package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_UnmarshallQueryResponse(t *testing.T) {
	body := []byte("{\"status\": \"test\"}")

	queryResponse, err := unmarshallQueryResponse(body)
	require.NoError(t, err)

	assert.Equal(t, "test", queryResponse.Status)
}
