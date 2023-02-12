package dbRepo

import (
	"context"
	"errors"
	"github.com/deenikarim/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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

//GetRoomByID go and gets a room by ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	//defined query for getting room by ID

	query := `
		select id, room_name, created_at, updated_at from rooms where id = $1`

	//execute the defined query //use when expecting only one row
	row := m.DB.QueryRowContext(ctx, query, id) //get the rows and then scan the row
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}

//TODO: CREATING DATABASE FUNCTIONS FOR AUTHENTICATION

//GetUserByID it goes and get a user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User

	//defined query for getting user by ID
	query := `
		select id, first_name, last_name, email, password, access_level,created_at, updated_at 
		from users where id = $1`

	//because we know we are getting only one row, we can use queryRowContext
	//execute the defined query //use when expecting only one row
	row := m.DB.QueryRowContext(ctx, query, id) //get the rows and then scan the row

	//copies the columns from the matched row into the values pointed to. todo: scan to the variable called u
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		//if there is an error
		return u, err
	}
	//otherwise no error
	return u, nil
}

//UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//defined query for updating users
	query := `
		update users set first_name = $1, last_name = $2, email = $3, access_level=$4, updated_at =$5
		`
	//execute query without returning any rows
	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.UpdatedAt,
	)

	if err != nil {
		//if there is an error
		return err
	}
	//otherwise no error
	return nil
}

//Authenticate performs the authentication of users
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//todo: store the information we get from the database
	//create a variable that will the ID of the authenticated user if things returned as they should
	//in other words, if they typed in the right password thus they have successfully been authenticated
	var id int
	//will hold the hash password for the authenticated user
	var hashedPassword string

	//todo: now we want to query the database, to see if we can find a user and store the information returned into the variables
	query := `
		select id, password from users where email = $1`

	////execute the defined query //use when expecting only one row which is email
	row := m.DB.QueryRowContext(ctx, query, email) //substitute query with empty value

	//scan the information received into some variable
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}
	//otherwise, at this point, if they have entered valid password and that matches an email in our database

	//todo: need to compare their passwords against our passwords [to compare their hash that we grab from
	//  database against a hash created from the password that the user typed in the form]
	/*TIP: all what are we doing is say, hey here is the hashedPassword that we pulled out of the database
	does this hash match the hash you are testing by running it against testPassword which is whatever
	the user typed in the form
	*/
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	//if they did match then all is good and were successful and can continue. Otherwise, not successful
	if err == bcrypt.ErrMismatchedHashAndPassword {
		//if there is error when trying to do comparison between hashes and passwords thus they do not match
		//want to do something else or return something
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		// if there is error and is something else other than mismatched
		return 0, "", err
	}

	//todo: if we get pass all of the above, then we are ready to return the necessary information because
	// the user can now be logged in
	return id, hashedPassword, nil
}
