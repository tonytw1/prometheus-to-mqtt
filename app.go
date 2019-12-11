package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkanos/gonfig"

	"github.com/tonytw1/prometheus-to-mqtt/domain"
	"github.com/tonytw1/prometheus-to-mqtt/prometheus"
)

type Configuration struct {
	PrometheusUrl string
	Jobs          []string
	MqttUrl       string
	MqttTopic     string
}

var client = http.Client{
	Timeout: time.Second * 5,
}

func main() {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	prometheusUrl := configuration.PrometheusUrl
	jobs := configuration.Jobs
	mqttURL := configuration.MqttUrl
	topic := configuration.MqttTopic

	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(mqttURL)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	defer c.Disconnect(250)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		// Publish metrics for each configured job
		for _, job := range jobs {
			// Metrics
			vectors, err := getMetrics(prometheusUrl, job)
			if err != nil {
				log.Print("Error getting metrics", err)
				continue
			}

			for _, instanceValue := range vectors {
				name := instanceValue.Metric["__name__"]
				value := instanceValue.Value[1].(string)

				message := formatMessage(job, name.(string), value)
				publish(c, topic, message)
			}

		}

		// Publish state of alerts - TODO are rules selectable by job?
		for _, rule := range getRules(prometheusUrl) {
			isAlertingRule := rule.Type == "alerting"
			if !isAlertingRule {
				continue
			}

			alertState := "0"
			for _, alert := range rule.Alerts {
				if alert.State == "firing" {
					alertState = "1"
					break
				}
			}
			publish(c, topic, formatMessage("", rule.Name, alertState)) // TODO job
		}

		time.Sleep(10 * time.Second)
	}
}

func formatMessage(job string, name string, value string) string {
	return job + "_" + name + ":" + value // TODO make safe
}

func getMetrics(prometheusUrl string, job string) ([]domain.InstantVector, error) {
	queryResponse, err := query(prometheusUrl, job)
	if err != nil {
		log.Print("Error while fetching", err)
		return nil, err
	}

	return prometheus.ExtractMetricsFromQueryResponse(queryResponse)
}

func getRules(prometheusUrl string) []domain.Rule {
	return []domain.Rule{}
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

	queryResponse, err := prometheus.UnmarshallQueryResponse(body)
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

func publish(c mqtt.Client, topic string, message string) {
	token := c.Publish(topic, 0, false, message)
	token.Wait()
}
