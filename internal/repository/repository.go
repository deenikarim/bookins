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
}
