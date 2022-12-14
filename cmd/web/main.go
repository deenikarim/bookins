package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/handlers"
	"github.com/deenikarim/bookings/internal/helpers"
	"github.com/deenikarim/bookings/internal/models"
	"github.com/deenikarim/bookings/internal/renders"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

/*instantiation our appConfig struct type and making universally accessible*/
var app config.AppConfig

//session making the sessionManager accessible to other packages which is the entire main package
var session *scs.SessionManager

//infoLog information about an error
var infoLog *log.Logger

//errorLog details about an error
var errorLog *log.Logger

func main() {

	err := run()
	if err != nil {
		log.Fatal(err) //write to the terminal and stop the application
	}

	//writing to the console
	fmt.Println(fmt.Sprintf("starting application on: %s", portNumber))

	//IMPROVED ROUTERS MODULE 1.1
	//how to use the routes() function created
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	//now starting the actual server
	err = srv.ListenAndServe()
	log.Fatal(err)

}

//run is for testing purposes
func run() error {
	//PART-3: SETTING UP APPLICATION WIDE CONFIGURATION AND NOW ; GO TO RENDER PACKAGE AND GET IT THERE
	/******************************************************************************************/

	//change this to true when in production, pull from appConfig
	app.InProduction = false

	//create the information log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog //store infoLog into appConfig

	//create the Error log
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog //store infoLog into appConfig

	//You need to tell golang: what am I going to put into the session(SESSION PART)
	gob.Register(models.Reservation{})

	//SESSIONS MODULE:5; SETTING UP A NEW SESSION MANAGER
	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //in production set to "true"

	//get the template cache into main.go
	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	//store template cache in struct type and session in the session struct type
	app.Session = session
	app.TemplateCache = tc
	app.UseCache = false
	/**/

	/*******************************************************************************************/

	/*from handler package ** INITIALIZE RENDER ***/
	/**/
	// NewRepo: calling the NewRepo in the main.go
	//getting our NewRepo function, outcome: create the repository variable
	repo := handlers.NewRepo(&app) //argument:: referencing to "app" to have access to appConfig struct type
	//after the repository variable is created, pass it back to NewHandlers
	handlers.NewHandlers(repo) //NOW:: to make changes to handlers function in order to have access to repository
	/**/

	/*from render package ** INITIALIZE RENDER ***/
	/**/
	//calling the NewTemplates function in the main.go ** INITIALIZE RENDER **
	renders.NewTemplates(&app)
	/**/

	//setting up the app variable and populating when the run() function is called
	helpers.NewHelpers(&app)

	return nil
}
