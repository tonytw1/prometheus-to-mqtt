package main

import (
	"encoding/json"
	"github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	prometheusUrl := "http://localhost"
	mqttUrl := "tcp://localhost:1883"
	mqttTopic := "test"

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(mqttUrl)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	queryResponse, err := query(prometheusUrl)
	if err != nil {
		log.Fatal(err)
	}

	if queryResponse.status != "success" {
		log.Fatal("Query response not successful")
	}
	if queryResponse.data.resultType != "vector" {
		log.Fatal("Expected result type vector")
	}

	data  := queryResponse.data
	var result = data.result.([]InstantVector)

	for _, instanceValue := range result {
		name := instanceValue.metric["__name__"]
		value := instanceValue.value[1].(string)

		message := name + ":" + value
		publish(c, mqttTopic, message)
	}

	c.Disconnect(250)
	println("Done")
}

func query(prometheusUrl string) (*QueryResponse, error) {
	queryUrl := prometheusUrl + "/api/v1/query"

	body, err := httpFetch(queryUrl)
	if err != nil {
		return nil, err
	}

	queryResponse := QueryResponse{}
	err = json.Unmarshal(body, &queryResponse)
	if err != nil {
		return nil, err
	}
	return &queryResponse, nil
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
