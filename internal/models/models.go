package models

import "time"

/*Reservation holds a reservation data structure (working server side validation)
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
*/

//User is the users model(describe our user)
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Room is the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//RoomRestriction is the room restrictions model
type RoomRestriction struct {
	ID                int
	StartDate         time.Time
	EndDate           time.Time
	RoomID            int
	ReservationID     int
	RestrictionTypeID int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Room              Room        //because roomID, it might include room information
	Reservation       Reservation //because reservationID, it might include reservation information
	RestrictionTypes  RestrictionTypes
}

//RestrictionTypes is the  restrictions types model
type RestrictionTypes struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservation is the reservation model
type Reservation struct {
	ID        int //from id to update_at needs to save in the database
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room //included all the room information here
}
