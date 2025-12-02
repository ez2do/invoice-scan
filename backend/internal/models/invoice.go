package models

type KeyValuePair struct {
	Key        string   `json:"key"`
	Value      string   `json:"value"`
	Confidence *float64 `json:"confidence,omitempty"`
}

type TableData struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

type InvoiceData struct {
	KeyValuePairs []KeyValuePair `json:"keyValuePairs"`
	Table         *TableData     `json:"table,omitempty"`
	Summary       []KeyValuePair `json:"summary"`
	Confidence    *float64       `json:"confidence,omitempty"`
}
