package prometheus

import (
	"encoding/json"
	"github.com/tonytw1/prometheus-to-mqtt/domain"
)

func UnmarshallQueryResponse(body []byte) (*domain.QueryResponse, error) {
	queryResponse := &domain.QueryResponse{}
	err := json.Unmarshal(body, queryResponse)
	if err != nil {
		return nil, err
	}
	return queryResponse, nil
}
