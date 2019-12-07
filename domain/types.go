package domain

type QueryData struct {
	ResultType string
	Result     interface{}
}

type QueryResponse struct {
	Status string
	Data   QueryData
}

type InstantVector struct {
	Metric map[string]interface{}
	Value  []interface{}
}

type Metric struct {
	Name string
}
