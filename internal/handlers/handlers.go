package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/deenikarim/bookings/internal/config"
	"github.com/deenikarim/bookings/internal/driver"
	"github.com/deenikarim/bookings/internal/forms"
	"github.com/deenikarim/bookings/internal/helpers"
	"github.com/deenikarim/bookings/internal/models"
	"github.com/deenikarim/bookings/internal/renders"
	"github.com/deenikarim/bookings/internal/repository"
	"github.com/deenikarim/bookings/internal/repository/dbRepo"
	"log"
	"net/http"
	"strconv"
	"strings"
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

//NewTestRepo enables us run tests on handlers package
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a, //populate the struct type created so everything in appConfig can be access by repository
		DB:  dbRepo.NewPostgresTestingRepo(a),
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
	renders.Template(w, r, "home.page.gohtml", &models.TemplateData{})
	m.DB.AllUsers()
}

//About create the about page handler function
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "about.page.gohtml", &models.TemplateData{})

}

//Contact create the Contact us page handler function
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "contact.page.gohtml", &models.TemplateData{})

}

//Generals create the generals page handler function
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "generals.page.gohtml", &models.TemplateData{})

}

//Majors create the Majors page handler function
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "majors.page.gohtml", &models.TemplateData{})

}

//SearchAvailability create the SearchAvailability page handler function
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})

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
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//calling the database function that does search for availability for all rooms
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room available!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	renders.Template(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})

	//w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))

}

//jsonResponse create struct for JSON request(link to CheckAvailabilityJSON function)
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	//modify the JSON to hand back roomID, start and end date
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//CheckAvailabilityJSON it handle request for checking availability and send JSON response
func (m *Repository) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	//calling the database function that allows searching for availability by room for a given date range
	//available variable is set to true if things are available and false if things are unavailable
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	//build request JSON request
	//INSTANTIATE THE STRUCT TYPE AND POPULATE IT
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
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
	//var emptyReservation models.Reservation

	// *******getting the reservation from the session***********
	//reservation at this point have some useful information in there
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) //type assertion
	if !ok {
		//improving errors handling by giving useful feedback
		m.App.Session.Put(r.Context(), "error", "can't get reservation from the session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	} //now have information like start and end dates and roomID

	//now we have to look up the room and stores its name in the correct location in our reservation
	// room contains room name, ID createdAt, and updatedAt
	room, err := m.DB.GetRoomByID(res.RoomID) //everything about the room is stored in this room variable
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//populating only the roomName field in the Room model
	res.Room.RoomName = room.RoomName

	//at this point it would make sense to store our reservation back into the session since we have already updated
	//our reservation model *****************************************
	m.App.Session.Put(r.Context(), "reservation", res) //contains start and end date, room name and room ID

	//take time.time and cast it back into a string in the format we expect
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	//how to get the formatted dates to the template
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	//AT THIS POINT

	//calling the renderTemplate function inside the handler function to render page to the browser
	renders.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		//send data to the template
		Form: forms.New(nil), //this includes or create  an empty form
		Data: data,           //adding an empty reservation input form
		//now passing the stringMap to the template
		StringMap: stringMap,
	})
}

//PostReservation handlers the posting of the reservation form (working server side validation)
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	/*/PULLING THE UPDATED VERSION OF THE RESERVATION MODEL FROM THE SESSION
	//START AND END DATE, ROOM ID, ROOM NAME
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) //type assert
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from the session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}*/

	//parsing the form data
	err := r.ParseForm()
	if err != nil {
		//now when something were wrong while parsing the form then this code will write to the terminal
		//the detailed error message
		m.App.Session.Put(r.Context(), "error", "can't parsed the form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "can't parsed the start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parsed the end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//convert output room id as a string or whatever into an integer
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return //if we don't add the RETURN, it will not stop executing at this point but i want to stop
	}

	/*/let's just update our reservation by adding first and last names, email and phone number
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")*/

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

		//http.Error(w, "my own error message", http.StatusSeeOther)
		//calling the renderTemplate function inside the handler function to render page to the browser
		renders.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			//send data to the template
			//Form: forms.New(nil), //this includes an empty form
			Form: form,
			Data: data, //passing the reservation: NOW have the information the user entered
		})
		return
	}

	//insert our reservation information to the database(save it to our database)
	//newReservationID purposes: use the reservation along with other information to build a room restriction
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	/*AFTER INSERTING THE RESERVATION WE NEED TO INSERT THE ROOM RESERVATION RESTRICTIONS BEFORE REDIRECT THEM*/
	//so now we are going to use the newReservationID along with all the required information to build up
	//a room restriction model and send it back to a function called InsertRoomRestriction()
	restriction := models.RoomRestriction{
		StartDate:         reservation.StartDate,
		EndDate:           reservation.EndDate,
		RoomID:            reservation.RoomID,
		ReservationID:     newReservationID,
		RestrictionTypeID: 1,
	}

	//write off our restriction information to the database(save it to our database)
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction into the database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//TODO: SEND NOTIFICATION THROUGH EMAIL TO GUEST
	//build the content of the email message
	htmlMessage := fmt.Sprintf(`
		<strong> Reservation Confirmation</strong><br>
		Dear %s: <br>
		This is to confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg := models.MailData{
		To:       reservation.Email,
		From:     "mike@example.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html", // will the name of the template inside the email-template folder
	}
	//send the message variable to the channel
	m.App.MailChan <- msg

	//TODO: SEND NOTIFICATION THROUGH EMAIL TO PROPERTY OWNER
	//build the content of the email message
	htmlMessage = fmt.Sprintf(`
		<strong> Reservation Notification</strong><br>
		Dear %s, <br>
		A reservation has been made for %s from %s to %s.
	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg = models.MailData{
		To:      "here@example.com",
		From:    reservation.Email,
		Subject: "Reservation Notification",
		Content: htmlMessage,
	}
	//send the message variable to the channel
	m.App.MailChan <- msg

	//todo:taking users to reservation summary page
	//how to put something into a session(now putting the reservation or whatever they have entered into a session)
	//NOW LETS PUT ON RESERVATION BACK TO THE SESSION
	m.App.Session.Put(r.Context(), "reservation", reservation)

	//redirect our user to reservation summary page(now when the user hit the submit button, goes to reservation summary page)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//ReservationSummary handles the posted form information(displaying a confirmation page)
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	//how to get something out of a session (reservation out of the session)
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

	//PUTTING THE DATES IN A FORMAT WE CAN DISPLAY
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data, //passing the reservation pulled out of the session to the template
		StringMap: stringMap,
	})

}

//ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	/*/TIP: Chi router has some helper methods that will allow us to that(to get the ID part of the URL)
	//get the ID part of the URL when a user clicks on any room link
	roomID, err := strconv.Atoi(chi.URLParam(r, "id")) //key is what we called it in the route file
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	} // now we have the room ID */

	// todo:That convenience function offered by chi, chi.URLPara(r, "id") is really, really hard to test
	//  In truth, we don't even need to use it, since we can parse the URL and find the id on our own, using this code:

	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	} //todo: In the code above, we use the strings package from the standard library to "explode" the URL into a slice of strings.
	//       Then, we grab the third element of that slice (position 2, since slices start counting from 0), and parse that into an int.

	//WHAT TO DO: need to get that reservation with only start and end dates that we store in the session
	// ---need to populate or update roomID field which has a zero value So let's stick our roomID variable from
	// ---above to the roomID field in our reservation model and then put it back in the session

	// *******getting the reservation from the session***********
	// Get returns the value for a given key from the session data. The return value has the type interface{}
	// so will usually need to be type asserted before you can use it

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) //type assertion
	if !ok {
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	} //now have the reservation stored in the session

	//update the roomID field with our variable from above roomID
	res.RoomID = roomID

	//then put it back in the session
	m.App.Session.Put(r.Context(), "reservation", res)

	//once they have clicked on a room, we want to display make reservation page
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//BookRoom take URL parameters then build a session variable and take users to make reservation page from
//General and Majors pages
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	//grabbing the values from the URL and are stored in parameters called
	//id, s, and e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parsed the start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parsed the end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//our session variable that is going to be put into the sessions
	//get our room name, start date and end date, roomID
	var res models.Reservation

	//get the room name and other details from the database
	rooms, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get the roomID from the database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = rooms.RoomName

	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	//at this point, we have the session variable that contains our roomID, roomName, startDate, endDate which is
	//needed in order to be able to go to make reservation.page

	//putting the session variable back into the session
	m.App.Session.Put(r.Context(), "reservation", res)

	//now, we want to take the users to make reservation page
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

//ShowLogin for displaying the login form for users to fill out to login in
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	//render login form
	renders.Template(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil), //pass it with an empty form
	})
}

//PostShowLogin handles the posting of login data and logging a user in.
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//well there are only two things coming out of the login form

	//todo step-1 call the RenewToken() method when going a log in or log out in order to prevent session fixation
	//  attacks
	_ = m.App.Session.RenewToken(r.Context())

	//todo step-2 parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//todo step-3 get the two things that come out of the login form
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	//todo step-4 making sure that our form has all the necessary parameters
	//creating a form object
	form := forms.New(r.PostForm) //contains all those url values and the associate data

	//passing it through validation
	form.Required("email", "password")
	//IsEmail checks for valid email address
	form.IsEmail("email")

	//if the form actually have errors or problems on it or in other words if it is not valid, then take the
	// user back
	if !form.Validate() {

		//we want to take the user back to the login screen if things didn't work the way we expected or
		// if the email and password are missing or don't match what we expect, then take the user back

		//calling the renderTemplate function inside the handler function to render page to the browser
		renders.Template(w, r, "login.page.gohtml", &models.TemplateData{
			//send data to the template
			Form: form, // pass it our form info
		})
		return
	}

	//todo: step-5 if it pass validation, then we want to authenticate the user
	//calling the authenticate() method we created in postgres.go and get the data entered from the login page and
	// use it for authentication
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	} //otherwise they have authenticated so can now log in

	//todo: step-6 if they succeeded from authenticating then we need to log the user in
	//we are going to that by storing the id returned from authentication in the session
	m.App.Session.Put(r.Context(), "user_id", id)

	//todo step-7 we need to take the user somewhere else after passing authentication process with an okay status
	m.App.Session.Put(r.Context(), "flash", "logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//LogOut handles the logging out of users
func (m *Repository) LogOut(w http.ResponseWriter, r *http.Request) {
	//TODO: how to log a user out of a system
	//so all we need to is to make sure that session variable called user_id does not exist
	//SIMPLE SOLUTION: IS TO JUST DESTROY THE SESSION
	_ = m.App.Session.Destroy(r.Context())
	//renew session token update the session data to have a new token while maintaining the current session data
	_ = m.App.Session.RenewToken(r.Context())

	//then we need to redirect them to somewhere else otherwise they're just going to see white screen
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

//AdminDashboard shows pages that are only available to logged-in users
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "admin-dashboard.page.gohtml", &models.TemplateData{})
}

// AdminAllReservations shows all reservations inu admin tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "admin-all-dashboard.page.gohtml", &models.TemplateData{})

}

// AdminNewReservations shows all new reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "admin-new-dashboard.page.gohtml", &models.TemplateData{})

}

// AdminReservationsCalendar displays the reservation calendar
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "admin-reservations-calendar.page.gohtml", &models.TemplateData{})

}

// AdminShowReservation shows the reservation in the admin tool
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {

}

// AdminPostShowReservation posts a reservation
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {

}

// AdminProcessReservation  marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

}

// AdminDeleteReservation deletes a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

}

// AdminPostReservationsCalendar handles post of reservation calendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {

}
