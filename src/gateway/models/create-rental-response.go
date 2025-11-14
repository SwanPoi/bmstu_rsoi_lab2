package models

type CreateRentalResponse struct {
	RentalUID string    				`json:"rental_uid"`
    CarUID    string    				`json:"car_uid"`
    DateFrom  string 					`json:"date_from"`
    DateTo    string 					`json:"date_to"`
    Status    string    				`json:"status"`
	Payment   PaymentCreationResponse	`json:"payment"`
}