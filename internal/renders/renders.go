package renders

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

//functions: holds all the functions we want to make available to our templates
var functions = template.FuncMap{
	"humanDate": HumanDate,
}

//pathToTemplate prevent hard coding of path
var pathToTemplate = "./templates"

/*PART-3: GETTING THE appCONFIG into the render package */
/**************************************************************************/

//app hold the type AppConfig (a pointer to AppConfig)
var app *config.AppConfig

//NewRenderer set the config settings for the templates
func NewRenderer(a *config.AppConfig) {
	app = a
}

//HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	//take the time that we passed and send it back in particular format
	return t.Format("2006-01-02")
}

/*************************************************** END ****************************************/

/**/
//PART-6: SHARING DEFAULT DATA

// AddDefaultData add data that should be available to every single page
//From
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	//populate the fields in template data(will automatically be populated when I am rendering page)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	//adding CSRF protection and using noSurf csrf package
	td.CSRFToken = nosurf.Token(r) //store CSRF token in the field of CSRFToken

	//todo determine whether or not a user is logged in
	// tip: if you recall, when a user is logged in, we setup session variable called user_id
	//so if the user id exists, then is authenticated
	if app.Session.Exists(r.Context(), "user_id") {
		//calling the member IsAuthenticated from the templateData.go
		//so if the value is greater than zero, then the user is logged in but if it is equal to zero then
		//the user is not logged in
		td.IsAuthenticated = 1 //default value for int is zero
	}

	return td // what it is doing now is taking the templateData and just returning
}

/**/

/********************* PART-2: RENDERING TEMPLATES TO THE BROWSER WINDOW **************************/

//Template a function for rendering templates
//what the function does is that it take a respondWriter and the name of a template you want to parse and read
// it to the browser
func Template(w http.ResponseWriter, r *http.Request, gohtml string, td *models.TemplateData) error {

	// get the template cache
	/*tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}*/
	/*******if in development mode, not production do not use the template cache instead
	rebuild the template on every  request ****************/
	var tc map[string]*template.Template //get tc from the scope of the if block//hold template cache
	if app.UseCache {
		tc = app.TemplateCache //using the template cache from app wide configuration
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get the individual templates from the myCache variable
	t, ok := tc[gohtml]
	if !ok {
		//log.Fatal("can not fetch the individual template")
		return errors.New("can't get item from the cache")
	}

	//store the template into a buffer
	buff := new(bytes.Buffer)

	//here is where to call AddDefaultData function before the execute function
	td = AddDefaultData(td, r)

	_ = t.Execute(buff, td)

	_, err := buff.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to the browser", err)
		return err
	}
	//if there is no error it returns nil
	return nil
}

/***************************************** END ***************************************************/

//******************* PART-1: CREATION OF TEMPLATE CACHE ******************************************

//CreateTemplateCache create template cache as a map that will hold all the templates
func CreateTemplateCache() (map[string]*template.Template, error) {
	//1.1: the variable myCache will hold all the template, thus it create it at the start of the app
	myCache := map[string]*template.Template{} //produces a safe html document fragment

	//1.2: find all the necessary pages in the template folder
	//Glob function returns the names of all files matching a pattern or nil if there is no match files
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplate))
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
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		//if there is at least one thing that it finds, then the length of matches will be greater than zero
		//if the above code can find any file that ends with .layout.html, then want to do something with it
		if len(matches) > 0 {
			//if it is greater than 0 or finds a file with that extension, what do I do with it
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*layout.gohtml", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		//adding the template set and the variable myCache
		myCache[name] = ts
	}
	return myCache, nil
}

/******************************************END*********************************************************/
