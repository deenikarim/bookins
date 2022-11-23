package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/deenikarim/bookings/pkg/config"
	"github.com/deenikarim/bookings/pkg/handlers"
	"github.com/deenikarim/bookings/pkg/renders"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

/*instantiation our appConfig struct type and making universally accessible*/
var app config.AppConfig

//session making the sessionManager accessible to other packages which is the entire main package
var session *scs.SessionManager

func main() {

	//PART-3: SETTING UP APPLICATION WIDE CONFIGURATION AND NOW ; GO TO RENDER PACKAGE AND GET IT THERE
	/******************************************************************************************/

	//change this to true when in production, pull from appConfig
	app.InProduction = false

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
	}
	//store template cache in struct type and session in the session struct type
	app.Session = session
	app.TemplateCache = tc
	app.UseCache = false
	/**/

	/*******************************************************************************************/

	/*from handler package*/
	/**/
	// NewRepo: calling the NewRepo in the main.go
	//getting our NewRepo function, outcome: create the repository variable
	repo := handlers.NewRepo(&app) //argument:: referencing to "app" to have access to appConfig struct type
	//after the repository variable is created, pass it back to NewHandlers
	handlers.NewHandlers(repo) //NOW:: to make changes to handlers function in order to have access to repository
	/**/

	/*from render package*/
	/**/
	//calling the NewTemplates function in the main.go
	renders.NewTemplates(&app)
	/**/

	//HandleFunc: registers the handler function for the given pattern(ROUTERS)
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	//http.HandleFunc("/reservation", handlers.Repo.Reservation)

	//writing to the console
	fmt.Println(fmt.Sprintf("starting application on: %s", portNumber))

	// listen on the TCP network address and then calls serve with handler to handle requests on incoming
	// connection(creating the web serve that listen to request)
	//_ = http.ListenAndServe(portNumber, nil)

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
