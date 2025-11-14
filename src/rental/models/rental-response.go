package models

import "time"

type RentalResponse struct {
    RentalUID string    `json:"rental_uid"`
    PaymentUID string   `json:"payment_uid"`
    CarUID    string    `json:"car_uid"`
    DateFrom  time.Time `json:"date_from"`
    DateTo    time.Time `json:"date_to"`
    Status    string    `json:"status"`
}