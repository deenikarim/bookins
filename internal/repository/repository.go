package repository

import (
	"github.com/deenikarim/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	//InsertReservation adding description of the function called InsertReservation
	InsertReservation(res models.Reservation) (int, error)

	//InsertRoomRestriction adding description of the function called InsertRoomRestriction
	InsertRoomRestriction(r models.RoomRestriction) error

	//SearchAvailabilityByDatesByRoomID search for availability for a given room
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)

	//SearchAvailabilityForAllRooms search for availability for all rooms
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	//GetRoomByID gets a room by ID
	GetRoomByID(id int) (models.Room, error)

	//GetUserByID returns a user by ID
	GetUserByID(id int) (models.User, error)

	//UpdateUser updates a user in the database
	UpdateUser(u *models.User) error

	//Authenticate performs the authentication
	Authenticate(email, testPassword string) (int, string, error)

	//AllReservations returns a list of all reservation
	AllReservations() ([]models.Reservation, error)

	//AllNewReservations displays all reservations with default value of 0
	AllNewReservations() ([]models.Reservation, error)

	//GetReservationByID returns one reservation by ID
	GetReservationByID(id int) (models.Reservation, error)
}
