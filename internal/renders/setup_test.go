package renders

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager

var testApp config.AppConfig

func TestMain(m *testing.M) {
	//GETTING THE SESSION INFORMATION FROM MAIN.GO

	//change this to true when in production, pull from appConfig
	testApp.InProduction = false

	//You need to tell golang: what am I going to put into the session(SESSION PART)
	gob.Register(models.Reservation{})

	//SESSIONS MODULE:5; SETTING UP A NEW SESSION MANAGER
	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //in production set to "true"

	//store template cache in struct type and session in the session struct type
	testApp.Session = session

	//it makes sure that "app" which is declared in render.go provides us with things we need to have
	app = &testApp

	os.Exit(m.Run())
	//code: this function get called before any of our test are run, it does whatever we tell it to do
	//in the body of the function, then just before it closes it run our test
}

//myWriter creates a struct type with no members
type myWriter struct {
}

//Header satisfied the header in ResponseWriter method
func (tw *myWriter) Header() http.Header {
	var h http.Header //create an empty header variable
	return h
}

//WriteHeader satisfied the WriteHeader in ResponseWriter method
func (tw *myWriter) WriteHeader(i int) {

}

//Write satisfy the write method in the ResponseWriter method
func (tw *myWriter) Write(b []byte) (int, error) {
	//NOW, in this case we can't just return a random int so we need to figure out the length of the slice of bytes
	length := len(b)

	return length, nil

}
