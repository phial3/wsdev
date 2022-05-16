package types

type ClientRequest struct {
	Action  string `json:"action"`
	Symbols string `json:"symbols,omitempty"`
}

type ClientResponse struct {
	Timestamp string  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float32 `json:"price"`
}

type Response struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []Data `json:"data"`
}

type Data struct {
	Timestamp string  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float32 `json:"lastPrice"`
}

type ResponseAuth struct {
	Success bool `json:"success"`
}
