package models

type PaymentInfo struct {
	PaymentUID string       `json:"payment_uid"`
    Status     string       `json:"status"`
    Price      int          `json:"price"`
}