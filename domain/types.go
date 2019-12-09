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

type Rule struct {
	Name string
}

type Group struct {
	Name  string
	Rules []Rule
}

type RulesData struct {
	Groups []Group
}

type RulesResponse struct {
	Status string
	Data   RulesData
}
