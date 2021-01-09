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

	response, err := ExtractMetricsFromQueryResponse(queryResponse)
	require.NoError(t, err)

	assert.Equal(t, 9, len(response))
}

func Test_UnmarshallRulesResponse(t *testing.T) {
	body, err := ioutil.ReadFile("rules.json")
	require.NoError(t, err)

	rulesResponse, err := UnmarshallRulesResponse(body)
	require.NoError(t, err)

	assert.Equal(t, "success", rulesResponse.Status)

	assert.Equal(t, 2, len(rulesResponse.Data.Groups))
	firstGroup := rulesResponse.Data.Groups[0]

	// The job a rule belongs to can be found on the name of the group the rule belongs in
	assert.Equal(t, "carbon-intensity", firstGroup.Name)

	assert.Equal(t, 2, len(firstGroup.Rules))
	alertingRule := firstGroup.Rules[1]
	assert.Equal(t, "AlwaysOnAlert", alertingRule.Name)
	assert.Equal(t, "alerting", alertingRule.Type)

	// Alerting rules have an alerts field
	assert.Equal(t, 1, len(alertingRule.Alerts))
	alert := alertingRule.Alerts[0]
	assert.Equal(t, "firing", alert.State, "We can tell if an alert is currently firing by looking at the state field")
}
