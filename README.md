# Prometheus to MQTT

We use Prometheus for monitoring and alerting but have several dashboard devices which consume metrics from a MQTT topic.

This application polls the Prometheus v1 API and exports metrics and alerts statues for a given set of jobs onto an MQTT topic.

Metrics are published as MQTT messages with the format:

```
[job name]/[metric name]:[integer value]
```

Alerts are published in the same format with a value of 0 or 1 to indicate if the alert is currently firing.

```
[job_name]/[alert name]:[true|false]
```
