package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonytw1/prometheus-to-mqtt/prometheus"
	"testing"
)

func Test_UnmarshallQueryResponse(t *testing.T) {
	body := []byte("{\"status\": \"test\"}")

	queryResponse, err := prometheus.UnmarshallQueryResponse(body)
	require.NoError(t, err)

	assert.Equal(t, "test", queryResponse.Status)
}
