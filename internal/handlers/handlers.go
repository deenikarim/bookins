package handlers

import (
	"encoding/json"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/driver"
	"github.com/deenikarim/bookings/internal/forms"
	"github.com/deenikarim/bookings/internal/helpers"
	"github.com/deenikarim/bookings/internal/models"
	"github.com/deenikarim/bookings/internal/renders"
	"github.com/deenikarim/bookings/internal/repository"
	"github.com/deenikarim/bookings/internal/repository/dbRepo"
	"net/http"
	"strconv"
	"time"
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
	DB repository.DatabaseRepo
}

//NewRepo create a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a, //populate the struct type created so everything in appConfig can be access by repository
		DB:  dbRepo.NewPostgresRepo(db.SQL, a),
		//reason: above code; db is not a repository but just a pointer to driver.DB so use the
		//NewPostgresRepo() instead
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
	renders.Template(w, r, "home.page.html", &models.TemplateData{})
	m.DB.AllUsers()
}

//About create the about page handler function
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "about.page.html", &models.TemplateData{})

}

//Contact create the Contact us page handler function
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "contact.page.html", &models.TemplateData{})

}

//Generals create the generals page handler function
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "generals.page.html", &models.TemplateData{})

}

//Majors create the Majors page handler function
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "majors.page.html", &models.TemplateData{})

}

//SearchAvailability create the SearchAvailability page handler function
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "search-availability.page.html", &models.TemplateData{})

}

//PostSearchAvailability create the handler for the POST request
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	//get the form data which is the input namespace ok of the request
	// this grabs the start and end from form post and stores it in start and end variables as strings because anything
	// you grab from a form post is a string
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	//convert our string from above to time.Time
	//// 2020-01-01 -- 01/02 03:04:05PM '06 -0700 format we expect our dates to be in
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//calling the database function that does search for availability for all rooms
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//let's now see what is the rooms variable by printing it to the screen ot terminal

	//when there is no availability there will be nothing in that slice(empty slice)
	if len(rooms) == 0 {
		//no availability by showing an error message
		//pass a value("error message") into our session
		m.App.Session.Put(r.Context(), "error", "no room available")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	} //otherwise, there is availability

	//making a map in order use it to store our rooms variable
	data := make(map[string]interface{})
	//store our rooms in there
	data["rooms"] = rooms

	//store some information entered in by the user like start and end dates before rendering the page
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		//RoomID:    0,leave it empty and populate that when they click on a room
	}
	//now store the res variable into the session so now we that necessary information available to us
	m.App.Session.Put(r.Context(), "reservation", res)

	//now need to render the template(choose-room template) and pass it data
	renders.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})

	//w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))

}

//jsonResponse create struct for JSON request(link to CheckAvailabilityJSON function)
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
		helpers.ServerError(w, err)
		return
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
	renders.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		//send data to the template
		Form: forms.New(nil), //this includes or create  an empty form
		Data: data,           //adding an empty reservation input form
	})
}

//PostReservation handlers the posting of the reservation form (working server side validation)
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	//parsing the form data
	err := r.ParseForm()
	//err = errors.New("this is something contains errors") ()DEMONSTRATING errors
	if err != nil {
		//now when something were wrong while parsing the form then this code will write to the terminal
		//the detailed error message
		helpers.ServerError(w, err)
		return
	}

	//converting our date in string format that our model expects
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	//convert our string from above to time.Time
	//// 2020-01-01 -- 01/02 03:04:05PM '06 -0700 format we expect our dates to be in
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//convert output room id as a string or whatever into an integer
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return //if we don't add the RETURN, it will not stop executing at this point but i want to stop
	}

	//initializes and populate the reservation struct Type
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"), //Get something from post request
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	//creating a form object
	form := forms.New(r.PostForm) //contains all those url values and the associate data

	//Has: it returns true or false and also add an error if first_name is empty
	//passing it through validation, for now it has one rule which is it has to have a first name
	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	//IsEmail checks for valid email address
	form.IsEmail("email")

	//if the form actually have errors or problems on it or in other words if it is not valid
	if !form.Validate() {
		//creating a variable to hold our reservation ::: creating a map to hold our string interface
		data := make(map[string]interface{})
		data["reservation"] = reservation //store the reservation variable from above in here

		//calling the renderTemplate function inside the handler function to render page to the browser
		renders.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			//send data to the template
			//Form: forms.New(nil), //this includes an empty form
			Form: form,
			Data: data, //passing the reservation: NOW have the information the user entered
		})
		return
	}

	//write off our reservation information to the database(save it to our database)
	//newReservationID purposes: use the reservation along with other information to build a room restriction
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	/*AFTER INSERTING THE RESERVATION WE NEED TO INSERT THE ROOM RESERVATION RESTRICTIONS BEFORE REDIRECT THEM*/

	//so now we are going to use the newReservationID along with all the required information to build up
	//a room restriction model and send it back to a function called InsertRoomRestriction()
	restriction := models.RoomRestriction{
		StartDate:         startDate,
		EndDate:           endDate,
		RoomID:            roomID,
		ReservationID:     newReservationID,
		RestrictionTypeID: 1,
	}

	//write off our restriction information to the database(save it to our database)
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
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
		//log.Println("can not get item from session")
		m.App.ErrorLog.Println("can not get item from session")
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
	renders.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data, //passing the reservation pulled out of the session to the template
	})

}
