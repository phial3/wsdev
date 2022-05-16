package config

type StockMarket struct {
	Scheme    string
	Address   string
	Path      string
	ApiKey    string
	Expired   int
	Signature string
	Topic     string
}

type Gateway struct {
	Host string
	Port int
}
