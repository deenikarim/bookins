package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/models"
	"github.com/deenikarim/bookings/internal/renders"
	"log"
	"net/http"
)

/**************************PART-4: USING THE APP_CONFIG IN THE HANDLER PACKAGE********************/
/******************* ALLOWING THE HANDLERS TO HAVE ACCESS TO APP CONFIG***************************/

//Repo create instance to use the Repository struct type
// the repository used by the handlers
var Repo *Repository

//Repository is the repository type which is a struct
// allow us to swap component of our application with minimum changes required to the base code
type Repository struct {
	App *config.AppConfig //embed a struct in another struct
	//things to put in here example sharing the database connection pool
}

//NewRepo create a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a, //populate the struct type created so everything in appConfig can be access by repository
	}
}

//NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	//the function does not return anything but set the variable created which is "Repo"
	Repo = r
}

//***************************************** END ************************************************//

//********************* PART-5: CREATING HANDLER FUNCTIONS**************************************//

//Home create the home page handler function
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

//About create the about page handler function
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "about.page.html", &models.TemplateData{})

}

//Contact create the Contact us page handler function
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})

}

//Generals create the generals page handler function
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})

}

//Majors create the Majors page handler function
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})

}

//SearchAvailability create the SearchAvailability page handler function
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})

}

//PostSearchAvailability create the handler for the POST request
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	//get the form data which is the input namespace ok of the request
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	//renders.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{})

}

//creating struct for JSON request(link to CheckAvailabilityJSON function)//BELOW FUNCTION USE THESE TYPE
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

//CheckAvailabilityJSON it handle request for checking availability and send JSON response
func (m *Repository) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	//build request JSON request
	//INSTANTIATE THE STRUCT TYPE AND POPULATE IT
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}
	//Convert the resp variable into JSON to create a JSON file and send it back
	outJson, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	//adding a Header that tells webs browser that is receiving my response what kind of response I am sending
	w.Header().Set("Content-Type", "application/json")

	//send the information directly to the ResponseWriter
	w.Write(outJson)
}

//Reservation create the reservation page handler function
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "hello again, says by kareem"

	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		//send data to the template
		StringMap: stringMap,
	})
}

//*********************************************** END ********************************************//
