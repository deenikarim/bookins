package dbRepo

import (
	"database/sql"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/repository"
)

type postgresDBRepo struct {
	//will hold two things
	App *config.AppConfig
	DB  *sql.DB //we are not going to populate it with  anything, but it needs to exist, so we can call
	//a database function that dont actually have database behind them
}

//testDBRepo for testing purposes
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB //holds the database connection pool
}

// NewPostgresRepo
//return takes two parameters in order to populate them with values
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

//NewPostgresTestingRepo for testing purposes
func NewPostgresTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
		//DB:  conn, not going to populate database part
	}
}
