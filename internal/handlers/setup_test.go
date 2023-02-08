package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/models"
	"github.com/deenikarim/bookings/internal/renders"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

//app holds the AppConfig
var app config.AppConfig

//session holds our session object
var session *scs.SessionManager

var pathToTemplate = "./../../templates"

//Function variable for the func type
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	//change this to true when in production, pull from appConfig
	app.InProduction = false

	//create the information log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog //store infoLog into appConfig

	//create the Error log
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog //store infoLog into appConfig

	//what am I going to put into the session(SESSION PART)
	gob.Register(models.Reservation{})

	//SESSIONS MODULE:5; SETTING UP A NEW SESSION MANAGER
	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //in production set to "true"

	//get the template cache into main.go
	tc, err := CreateTestTemplateCache()
	//tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("can't create template cache")
		//return nil
	}
	//store template cache in struct type and session in the session struct type
	app.Session = session
	app.TemplateCache = tc
	app.UseCache = true
	/**/

	/*******************************************************************************************/

	/*from handler package ** INITIALIZE RENDER ***/
	/**/
	// NewRepo: calling the NewRepo in the main.go
	//getting our NewRepo function, outcome: create the repository variable
	repo := NewTestRepo(&app) //argument:: referencing to "app" to have access to appConfig struct type
	//after the repository variable is created, pass it back to NewHandlers
	NewHandlers(repo) //NOW:: to make changes to handlers function in order to have access to repository
	/**/

	/*from render package ** INITIALIZE RENDER ***/
	/**/
	//calling the NewTemplates function in the main.go ** INITIALIZE RENDER **
	renders.NewRenderer(&app)
	/**/

	os.Exit(m.Run()) //before it exits run the tests

}

//getRoute we need to have access to our routes otherwise we can't call the handlers at all
//since it's going to give us routes, it will hand back exactly what our route does
func getRoute() http.Handler {

	//NOW GETTING OUR ROUTES
	mux := chi.NewRouter()

	/**/ //installing middleware
	mux.Use(middleware.Recoverer)

	/**/ //how to use the middleware created
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	/**/ //here is where to set up the routes

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/search-availability", Repo.SearchAvailability)

	//route to handle the POST request method
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", Repo.CheckAvailabilityJSON)

	//fileServer: creating fileServer to enable static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	/**/
	return mux
}

//IT FIX THE MIDDLEWARE ERRORS

//NoSurf create CSRFToken: adds CSRF protection to all POST Request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next) //will create a new handler

	//need to set some values for it
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", //the way to refer to entire site for a cookie
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad create middleware that loads  and saves the sessions on every request to make web server state aware
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//CreateTestTemplateCache create template cache as a map that will hold all the templates
//REASON: do not want to call this function directly
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//1.1: the variable myCache will hold all the template, thus it create it at the start of the app
	myCache := map[string]*template.Template{} //produces a safe html document fragment

	//1.2: find all the necessary pages in the template folder
	//Glob function returns the names of all files matching a pattern or nil if there is no match files
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplate))
	//checking for error because Glob function also returns an error if it finds no files
	if err != nil {
		return myCache, err
	} //now we have our pages but have not done anything it yet

	//1.3: iteration: get all the page.html file
	//loop through that range, arrange for those pages and print out the name of the current page
	for _, page := range pages {
		//what to do by looping through
		//now have the info about the files

		//extracting the actually base name because what it is returning is the full path to files
		name := filepath.Base(page)

		//now with the actual names of the pages, lets create a template set
		//create a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		//New(name): allocate a new html template with the given name
		//ParseFiles(page): parse the named files and associate the resulting template with
		if err != nil {
			return myCache, err
		}

		// Agenda: find out, does a template matches any layout or should use a specific layout defined purposeful
		// for a particular template
		// BEGIN: check to see if something matches

		//code below: look for any file in the template folder that end (.layout) or checking for the existence
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		//if there is at least one thing that it finds, then the length of matches will be greater than zero
		//if the above code can find any file that ends with .layout.html, then want to do something with it
		if len(matches) > 0 {
			//if it is greater than 0 or finds a file with that extension, what do I do with it
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*layout.html", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		//adding the template set and the variable myCache
		myCache[name] = ts
	}
	return myCache, nil
}
