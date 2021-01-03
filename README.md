# Prometheus to MQTT

Polls Prometheus and exports metrics and alerts for a given set of jobs onto an MQTT topic.

We use Prometheus for monitoring and alerting but have several dashboard devices which 
consumer metrics from a MQTT topic.

Metrics are published as messages with the format:

```
[metric name]:[integer value]
```

Alerts are published in the same format with a value of 0 or 1 to indicate if the alert is 
currently firing.

```
[alert name]:[true|false]
```

