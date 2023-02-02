package dbRepo

import (
	"context"
	"github.com/deenikarim/bookings/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation writes a routine that is going to write a reservation to the database
//use: insert a reservation into the database (now have the means of storing the needed reservation info)
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	//3*time.Second: if this  transaction take longer than 3s then cancel it and everything will go
	//back to where it should be
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	//writing an query called stmt(inserting)
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date,
			room_id, created_at, updated_at) 
			values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id` //getting the id of the last inserted record

	//now how to execute that stmt?
	err := m.DB.QueryRowContext(ctx, stmt,
		//passing it all the arguments
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID) //scanning the old thing that is being returned into the variable newID

	//Exec() function returns a result and error
	if err != nil {
		return 0, err //0 because if there is error we dont care about what the ID is
	}

	return newID, nil
}

//InsertRoomRestriction inserts a new room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	//3*time.Second: if this  transaction take longer than 3s then cancel it and everything will go
	//back to where it should be
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservations_id, created_at, updated_at,
    		restriction_type_id)
			values ($1,$2,$3,$4,$5,$6,$7)`

	//now how to execute that stmt?
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionTypeID,
	)

	if err != nil {
		return err
	}

	return nil
}

//SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID but false if not availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	//3*time.Second: if this  transaction take longer than 3s then cancel it and everything will go
	//back to where it should be
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//for setting up return
	var numRow int

	//now the where clause limit our search to a given room
	query := `
			select
				count(id) 
			from
				room_restrictions
			where 
				room_id = $1 and $2 < end_date and $3 > start_date;`

	//execute the defined query //use when expecting only one row
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end) //get the rows and then scan the row
	err := row.Scan(&numRow)
	if err != nil {
		return false, err //if there is error return false and error
	}

	//otherwise, if numRows is 0, then there is availability so return true
	if numRow == 0 {
		return true, nil
	}
	//otherwise, if numRows is greater than 0 then no availability so return false
	return false, nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms if any, for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	//defined query for searching for availability for all rooms
	query := `
			select
				r.id, r.room_name 
			from
				rooms r 
			where 
				r.id not in (select rr.room_id from room_restrictions rr 
				where $1 < rr.end_date and $2 > rr.start_date);`

	//return rows
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	//the results received to be sent to all destination
	for rows.Next() {
		var room models.Room
		//where to scan the  rows returned to
		err := rows.Scan(
			//destination to variable room
			//only getting two things (ID And Room name)
			&room.ID,
			&room.RoomName,
		)
		//checking for errors
		if err != nil {
			return rooms, err
		}
		//needs to put it into our slice of names 'rooms'
		rooms = append(rooms, room)
	}

	//error checking again
	err = rows.Err()
	if err != nil {
		return rooms, err
	}

	return rooms, nil
}
