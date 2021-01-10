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

	var logConnectionLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Print("Connection lost")
	}

	var logReconnecting mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Print("Reconnecting")
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttURL)
	opts.SetOnConnectHandler(logConnection)
	opts.SetConnectionLostHandler(logConnectionLost)
	opts.SetReconnectingHandler(logReconnecting)
	opts.SetClientID("prometheus-to-mqtt")

	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG]  ", 0)

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
			//println("Getting metrics for job: " + job)
			vectors, err := prometheus.GetMetrics(prometheusUrl, job)
			if err != nil {
				log.Print("Error getting metrics", err)
				continue
			}

			//println("Got vectors: " + strconv.FormatInt(int64(len(vectors)), 10))
			for _, instanceValue := range vectors {
				name := instanceValue.Metric["__name__"]
				value := instanceValue.Value[1].(string)
				publish(c, topic, formatMessage(job, name.(string), value))
			}

		}

		// Publish state of alerts
		ruleGroups, err := prometheus.GetRuleGroups(prometheusUrl)
		if err != nil {
			log.Print("Error getting rules", err)
			continue
		}

		for _, group := range ruleGroups {
			//job := group.Name
			for _, rule := range group.Rules {
				isAlertingRule := rule.Type == "alerting"
				if !isAlertingRule {
					continue
				}
				//alertState := "false"
				for _, alert := range rule.Alerts {
					if alert.State == "firing" {
						//alertState = "true"
						break
					}
				}
				//publish(c, topic, formatMessage(job, rule.Name, alertState))
			}
		}

		time.Sleep(5 * time.Second)
	}

	log.Print("End")
}

func formatMessage(job string, name string, value string) string {
	return job + "_" + name + ":" + value // TODO make safe
}

func publish(c mqtt.Client, topic string, message string) {
	token := c.Publish(topic, 0, false, message)
	token.WaitTimeout(time.Second * 1)
}
