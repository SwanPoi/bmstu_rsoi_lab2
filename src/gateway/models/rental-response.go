package models

type RentalResponse struct {
	RentalUID string    		`json:"rental_uid"`
    DateFrom  string 			`json:"date_from"`
    DateTo    string 			`json:"date_to"`
    Status    string    		`json:"status"`
	Car		  CarInfo 			`json:"car"`
	Payment   PaymentInfo		`json:"payment"`
}