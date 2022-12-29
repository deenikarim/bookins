package driver

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

//DB holds the connection pool
type DB struct {
	SQL *sql.DB
}

//dbConn initializes the struct type of DB
var dbConn = &DB{}

//native of our connection pool
//maxOpenDbConn set the maximum number of connections pool so never have more than 10 connections to the database open
//at any given time
const maxOpenDbConn = 10

//maxIdleDbConn means how many connections can remain in the connection pool but remain idle
const maxIdleDbConn = 5

//maxDbLifetime defines the maximum lifetime of a database connection pool
const maxDbLifetime = 5 * time.Minute

//ConnectSQL creates a database connection pool for postgres by calling another function
func ConnectSQL(dsn string) (*DB, error) {
	//calling the database connection pool
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err) //just die because we can't connect to the database so no further actions
	}

	//if everything goes well from the above code "d" variable will have some data or values
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	//store the database connection pool to the member SQL
	dbConn.SQL = d

	//let's test the connection pool again
	err = testDB(d)
	if err != nil {
		return nil, err
	}

	//return our database connection pool and nil if there is no error
	return dbConn, nil
}

//testDB for pinging database and establishing active connection pool
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database and connects to the database(our connection pool)
func NewDatabase(dsn string) (*sql.DB, error) {
	//connect to database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		//return nil and error because we don't have database connection pool
		return nil, err
	}
	//otherwise keep going
	//test our connection pool
	err = db.Ping()
	if err != nil {
		//return nil and error because we couldn't test our connection pool and establish active connection
		return nil, err
	}
	//so if everything works properly then return nil
	return db, nil
}
