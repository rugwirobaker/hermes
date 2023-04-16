package pindo

type SendRequest struct {
	To     string `json:"to"`
	Text   string `json:"text"`
	Sender string `json:"sender"`
}

type SendResponse struct {
	Bonus            float64 `json:"bonus"`
	ItemCount        int64   `json:"item_count"`
	Discount         float64 `json:"discount"`
	RemainingBalance string  `json:"remaining_balance"`
	SmsID            string  `json:"sms_id"`
	SelfURL          string  `json:"self_url"`
	ItemPrice        float64 `json:"item_price"`
	Status           string  `json:"status"`
	To               string  `json:"to"`
	TotalCost        float64 `json:"total_cost"`
}
