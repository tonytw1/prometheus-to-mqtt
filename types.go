package main

type QueryData struct {
	resultType string
	result 	interface{}
}

type QueryResponse struct {
	status	string
	data	QueryData
}

type InstantVector struct {
	metric map[string]string
	value []interface{}
}

type Metric struct {
	name string
	value float64
}

