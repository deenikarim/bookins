package dbRepo

import (
	"errors"
	"github.com/deenikarim/bookings/internal/models"
	"log"
	"time"
)

//FOR SOLVING THE MISSING METHODS FUNCTION BECAUSE WE NEED TO HAVE

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation writes a routine that is going to write a reservation to the database
//use: insert a reservation into the database (now have the means of storing the needed reservation info)
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//help test, attempt to insert a reservation into the database in case of any failure
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

//InsertRoomRestriction inserts a new room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	//help test, attempt to insert a reservation into the database in case of any failure
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

//SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID but false if not availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms if any, for a given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}
	if start == testDateToFail {
		return rooms, errors.New("some error")
	}
	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

//GetRoomByID go and gets a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	//help fail to getting of roomID
	if id > 2 {
		return room, errors.New("error because non-existing room")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	//TODO implement me
	//panic("implement me")
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u *models.User) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	//TODO implement me
	//panic("implement me")
	return 1, "", nil
}
