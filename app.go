package main

import (
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkanos/gonfig"

	"github.com/tonytw1/prometheus-to-mqtt/prometheus"
)

type Configuration struct {
	PrometheusUrl string
	Jobs          []string
	MqttUrl       string
	MqttTopic     string
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

	var logConnection mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Print("Connected")
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttURL)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetOnConnectHandler(logConnection)
	opts.SetCleanSession(true)
	opts.SetClientID("prometheus-to-mqtt")
	opts.SetAutoReconnect(true)

	println("Connecting to: ", mqttURL)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(250)

	for {
		println("Polling")
		// Publish metrics for each configured job
		for _, job := range jobs {
			// Metrics
			vectors, err := prometheus.GetMetrics(prometheusUrl, job)
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

		// Publish state of alerts
		// TODO are rules selectable by job?
		rules, err := prometheus.GetRules(prometheusUrl)
		if err != nil {
			log.Print("Error getting rules", err)
			continue
		}
		for _, rule := range rules {
			isAlertingRule := rule.Type == "alerting"
			if !isAlertingRule {
				continue
			}

			alertState := "false"
			for _, alert := range rule.Alerts {
				if alert.State == "firing" {
					alertState = "true"
					break
				}
			}
			publish(c, topic, formatMessage("", rule.Name, alertState)) // TODO job
		}

		time.Sleep(10 * time.Second)
	}

	log.Print("End")
}

func formatMessage(job string, name string, value string) string {
	return job + "_" + name + ":" + value // TODO make safe
}

func publish(c mqtt.Client, topic string, message string) {
	token := c.Publish(topic, 0, false, message)
	timeout := token.WaitTimeout(time.Second * 1)
	println(timeout)
}
