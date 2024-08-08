package models

type Client struct {
	ClientID string `json:"client_id"`
	Name     string `json:"name"`
}

type Deposit struct {
	DepositID string    `json:"deposit_id"`
	ClientID  string    `json:"client_id"`
	Nominal   float64   `json:"nominal"`
	CreatedAt string    `json:"created_at,omitempty"`
	Receipts  []Receipt `json:"receipts,omitempty"`
}

type Receipt struct {
	ReceiptID   string  `json:"receipt_id"`
	DepositID   string  `json:"deposit_id"`
	Amount      float64 `json:"amount"`
	CreatedAt   string  `json:"created_at,omitempty"`
	Pot         string  `json:"pot"`
	WrapperType string  `json:"wrapper_type"`
}
