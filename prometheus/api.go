package prometheus

import (
	"github.com/tonytw1/prometheus-to-mqtt/domain"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

func GetMetrics(prometheusUrl string, job string) ([]domain.InstantVector, error) {
	queryResponse, err := query(prometheusUrl, job)
	if err != nil {
		log.Print("Error while fetching", err)
		return nil, err
	}

	return ExtractMetricsFromQueryResponse(queryResponse)
}

func GetRules(prometheusUrl string) []domain.Rule {
	return []domain.Rule{} // TODO implement
}

func query(prometheusUrl string, job string) (*domain.QueryResponse, error) {
	queryUrl, err := url.Parse(prometheusUrl + "/api/v1/query")
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(queryUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	q := "{job=\"" + job + "\"}"
	values.Add("query", q)
	queryUrl.RawQuery = values.Encode()

	body, err := httpFetch(queryUrl.String())
	if err != nil {
		return nil, err
	}

	queryResponse, err := UnmarshallQueryResponse(body)
	if err != nil {
		return nil, err
	}
	return queryResponse, nil
}

func httpFetch(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	return body, nil
}
