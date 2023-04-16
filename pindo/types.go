package pindo

type SendRequest struct {
	To     string `json:"to"`
	Text   string `json:"text"`
	Sender string `json:"sender"`
}

type SendResponse struct {
	Bonus            float64 `json:"bonus"`
	Discount         float64 `json:"discount"`
	ItemCount        int64   `json:"item_count"`
	ItemPrice        float64 `json:"item_price"`
	RemainingBalance float64 `json:"remaining_balance"`
	SmsID            int     `json:"sms_id"`
	SelfURL          string  `json:"self_url"`
	Status           string  `json:"status"`
	To               string  `json:"to"`
	TotalCost        float64 `json:"total_cost"`
}
