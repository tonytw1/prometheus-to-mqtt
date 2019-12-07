package prometheus

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func Test_UnmarshallMetricsQueryResponse(t *testing.T) {
	body, err := ioutil.ReadFile("query.json")
	require.NoError(t, err)

	queryResponse, err := UnmarshallQueryResponse(body)
	require.NoError(t, err)

	assert.Equal(t, "success", queryResponse.Status)
}

func Test_ExtractMetricsFromQueryResponse(t *testing.T) {
	body, err := ioutil.ReadFile("query.json")
	require.NoError(t, err)

	queryResponse, err := UnmarshallQueryResponse(body)
	require.NoError(t, err)

	response := ExtractMetricsFromQueryResponse(queryResponse)

	assert.Equal(t, 9, len(response))
}

func Test_UnmarshallRulesResponse(t *testing.T) {
	body, err := ioutil.ReadFile("rules.json")
	require.NoError(t, err)

	rulesResponse, err := UnmarshallRulesResponse(body)
	require.NoError(t, err)

	assert.Equal(t, "success", rulesResponse.Status)
}
