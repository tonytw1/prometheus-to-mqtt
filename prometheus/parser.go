package prometheus

import (
	"encoding/json"
	"github.com/tonytw1/prometheus-to-mqtt/domain"
	"log"
)

func UnmarshallQueryResponse(body []byte) (*domain.QueryResponse, error) {
	queryResponse := &domain.QueryResponse{}
	err := json.Unmarshal(body, queryResponse)
	if err != nil {
		return nil, err
	}
	return queryResponse, nil
}

func ExtractMetricsFromQueryResponse(queryResponse *domain.QueryResponse) []domain.InstantVector {
	if queryResponse.Status != "success" {
		log.Fatal("Query response not successful")
	}
	if queryResponse.Data.ResultType != "vector" {
		log.Fatal("Expected result type vector")
	}

	data := queryResponse.Data
	result := data.Result.([]interface{})

	ivs := []domain.InstantVector{}
	for _, i := range result {
		j := i.(map[string]interface{})
		x := j["metric"].(map[string]interface{})
		y := j["value"].([]interface{})
		ivs = append(ivs, domain.InstantVector{Metric: x, Value: y})
	}

	return ivs
}

func UnmarshallRulesResponse(body []byte) (*domain.RulesResponse, error) {
	rulesResponse := &domain.RulesResponse{}
	err := json.Unmarshal(body, rulesResponse)
	if err != nil {
		return nil, err
	}
	return rulesResponse, nil
}
