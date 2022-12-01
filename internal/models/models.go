package models

//Reservation holds a reservation data structure (working server side validation)
type Reservation struct {
	FirstName          string
	LastName           string
	Email              string
	Phone              string
	Address            string
	AddressTwo         string
	City               string
	TermsAndConditions string
	State              string
}
