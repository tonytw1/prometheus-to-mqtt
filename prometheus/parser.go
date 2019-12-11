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

func ExtractMetricsFromQueryResponse(queryResponse *domain.QueryResponse) ([]domain.InstantVector, error) {
	if queryResponse.Status != "success" {
		log.Print("Query response not successful")
		return nil, nil
	}
	if queryResponse.Data.ResultType != "vector" {
		log.Print("Expected result type vector")
		return nil, nil
	}

	data := queryResponse.Data
	result := data.Result.([]interface{})

	ivs := []domain.InstantVector{}
	for _, i := range result {
		j := i.(map[string]interface{})
		ivs = append(ivs, domain.InstantVector{
			Metric: j["metric"].(map[string]interface{}),
			Value:  j["value"].([]interface{})})
	}

	return ivs, nil
}

func UnmarshallRulesResponse(body []byte) (*domain.RulesResponse, error) {
	rulesResponse := &domain.RulesResponse{}
	err := json.Unmarshal(body, rulesResponse)
	if err != nil {
		return nil, err
	}
	return rulesResponse, nil
}
