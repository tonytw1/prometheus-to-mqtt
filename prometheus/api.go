package prometheus

import (
	"errors"
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
		log.Print("Error while fetching metrics", err)
		return nil, err
	}

	return ExtractMetricsFromQueryResponse(queryResponse)
}

func GetRules(prometheusUrl string) ([]domain.Rule, error) {
	rulesResponse, err := fetchRules(prometheusUrl)
	if err != nil {
		log.Print("Error while fetching rules", err)
		return nil, err
	}

	if rulesResponse.Status != "success" {
		return nil, errors.New("Rules response status was not success")
	}

	var rules []domain.Rule
	for _, group := range rulesResponse.Data.Groups {
		for _, rule := range group.Rules {
			rules = append(rules, rule)
		}
	}
	return rules, nil
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

	q := "{job=\"" + job + "\"}" // TODO make safe
	values.Add("query", q)
	queryUrl.RawQuery = values.Encode()

	body, err := httpFetch(queryUrl)
	if err != nil {
		return nil, err
	}

	return UnmarshallQueryResponse(body)
}

func fetchRules(prometheusUrl string) (*domain.RulesResponse, error) {
	rulesUrl, err := url.Parse(prometheusUrl + "/api/v1/rules")
	if err != nil {
		return nil, err
	}

	body, err := httpFetch(rulesUrl)
	if err != nil {
		return nil, err
	}

	return UnmarshallRulesResponse(body)
}

func httpFetch(url *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}
