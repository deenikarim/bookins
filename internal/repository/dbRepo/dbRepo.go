package dbRepo

import (
	"database/sql"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/repository"
)

type postgresDBRepo struct {
	//will hold two things
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
