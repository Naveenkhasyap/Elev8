package tokensvc

type TokenData struct {
	Name        string `json:"name"`
	Ticker      string `json:"ticker"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Twitter     string `json:"twitter,omitempty"`
	Telegram    string `json:"telegram,omitempty"`
	Website     string `json:"website,omitempty"`
}
