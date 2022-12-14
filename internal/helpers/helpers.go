package helpers

import (
	"fmt"
	"github.com/deenikarim/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

//what to put in this helpers package

//app allows access to our config.AppConfig
var app *config.AppConfig

//our method of populating the variable declared with whatever is in our config.AppConfig

//NewHelpers set up our AppConfig for helpers(populating it with the contents of our config.AppConfig)
//we initialize this function from main.go
func NewHelpers(a *config.AppConfig) {
	//passing or equal to this function our config.AppConfig
	app = a
}

//ClientError handles client errors
func ClientError(w http.ResponseWriter, status int) {

	//writing to the InfoLog
	app.InfoLog.Println("client error with status of", status)

	//because it is a client error need to give some kind response
	//Error replies to the request with the specified error message and HTTP code.
	//The error message should be plain text
	http.Error(w, http.StatusText(status), status)
}

//ServerError handles something were wrong with the server
func ServerError(w http.ResponseWriter, err error) {

	//trace the nature of the error and will hold a string
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//writing to the terminal our errorLog
	app.ErrorLog.Println(trace)
	//give some feedback to the user which why we have the ResponseWriter
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
