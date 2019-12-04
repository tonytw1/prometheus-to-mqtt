package main

type QueryData struct {
	ResultType string
	Result     interface{}
}

type QueryResponse struct {
	Status string
	Data   QueryData
}

type InstantVector struct {
	metric map[string]interface{}
	value  []interface{}
}

type Metric struct {
	name string
}
