package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	PrometheusUrl string
	Job           string
	MqttUrl       string
	MqttTopic     string
}

func main() {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(configuration.MqttUrl)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		for _, instanceValue := range getMetrics(configuration.PrometheusUrl, configuration.Job) {
			name := instanceValue.metric["__name__"]
			value := instanceValue.value[1].(string)

			message := name.(string) + ":" + value
			publish(c, configuration.MqttTopic, message)
		}

		time.Sleep(10 * time.Second)
	}

	c.Disconnect(250)
}

func getMetrics(prometheusUrl string, job string) []InstantVector {
	queryResponse, err := query(prometheusUrl, job)
	if err != nil {
		log.Fatal(err)
	}

	if queryResponse.Status != "success" {
		log.Fatal("Query response not successful")
	}
	if queryResponse.Data.ResultType != "vector" {
		log.Fatal("Expected result type vector")
	}

	data := queryResponse.Data
	result := data.Result.([]interface{})

	ivs := []InstantVector{}
	for _, i := range result {
		j := i.(map[string]interface{})
		x := j["metric"].(map[string]interface{})
		y := j["value"].([]interface{})
		ivs = append(ivs, InstantVector{metric: x, value: y})
	}

	return ivs
}
func query(prometheusUrl string, job string) (*QueryResponse, error) {
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

	queryResponse, err := unmarshallQueryResponse(body)
	if err != nil {
		return nil, err
	}
	return queryResponse, nil
}

func unmarshallQueryResponse(body []byte) (*QueryResponse, error) {
	queryResponse := &QueryResponse{}
	err := json.Unmarshal(body, queryResponse)
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

	client := http.Client{
		Timeout: time.Second * 5,
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
