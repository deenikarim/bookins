package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/forms"
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
	//empty Reservation but it has all the necessary fields
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation //stores the empty reservation "TIP: needs to have the same name used to store the reservation in POST-RESERVATION"

	//calling the renderTemplate function inside the handler function to render page to the browser
	renders.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		//send data to the template
		Form: forms.New(nil), //this includes an empty form
		Data: data,           //adding an empty reservation input form
	})
}

//PostReservation handlers the posting of the reservation form (working server side validation)
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	//parsing the form data
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	//initializes and populate the reservation struct Type
	reservation := models.Reservation{
		FirstName:          r.Form.Get("first_name"), //Get something from post request
		LastName:           r.Form.Get("last_name"),
		Email:              r.Form.Get("email"),
		Phone:              r.Form.Get("phone"),
		Address:            r.Form.Get("address"),
		AddressTwo:         r.Form.Get("address_two"),
		City:               r.Form.Get("city"),
		TermsAndConditions: r.Form.Get("terms_and_conditions"),
		State:              r.Form.Get("state"),
	}

	//creating a form object
	form := forms.New(r.PostForm) //contains all those url values and the associate data

	//Has: it returns true or false and also add an error if first_name is empty
	//passing it through validation, for now it has one rule which is it has to have a first name
	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email", "phone", "address", "address_two", "city")
	form.MinLength("first_name", 3, r)
	//IsEmail checks for valid email address
	form.IsEmail("email")

	//if the form actually have errors or problems on it or in other words if it is not valid
	if !form.Validate() {
		//creating a variable to hold our reservation ::: creating a map to hold our string interface
		data := make(map[string]interface{})
		data["reservation"] = reservation //store the reservation variable from above in here

		//calling the renderTemplate function inside the handler function to render page to the browser
		renders.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			//send data to the template
			//Form: forms.New(nil), //this includes an empty form
			Form: form,
			Data: data, //passing the reservation: NOW have the information the user entered
		})
		return
	}

	//taking users to reservation summary page
	//how to put something into a session(now putting the reservation or whatever they have entered into a session)
	m.App.Session.Put(r.Context(), "reservation", reservation)

	//redirect our user to reservation summary page(now when the user hit the submit button, goes to reservation summary page)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//ReservationSummary handles the posted form information(displaying a confirmation page)
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	//how to get something out of a session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) //type assert
	if !ok {
		log.Println("can not get item from session")
		//pass a value("error message") into our session
		m.App.Session.Put(r.Context(), "error", "can not get item from session")
		//still return a blank page, so we now need to redirect them to somewhere else which is the HOME Page
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//taken out our session from reservation post handler
	m.App.Session.Remove(r.Context(), "reservation") //remove data from our reservation

	//after getting the reservation item from the session then I want to do something with it(pass the reservation pulled as a template data)
	data := make(map[string]interface{})
	data["reservation"] = reservation //store the reservation variable created above in our TemplateData under field "Data"

	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data, //passing the reservation pulled out of the session to the template
	})

}

//*********************************************** END ********************************************//
