package globe

type CurrencyInfoType struct {
	// FlagURL string `json:"flagURL,omitempty"`
	Name string `json:"name,omitempty"`
	Rate string `json:"rate"`
	// RateNm string `json:"ratenm"`
	Scur   string `json:"scur"`
	Status string `json:"status"`
	Symbol string `json:"symbol,omitempty"`
	Tcur   string `json:"tcur"`
	Update string `json:"update"`
	Hot    int    `json:"hot,omitempty"`
}

type CurrencyInfoMapType map[string]CurrencyInfoType

type Result struct {
	Status string `json:"status"`
	Scur   string `json:"scur"`
	Tcur   string `json:"tcur"`
	Ratenm string `json:"ratenm"`
	Rate   string `json:"rate"`
	Update string `json:"update"`
}
type Res struct {
	Success string `json:"success"`
	Result  Result `json:"result"`
}
