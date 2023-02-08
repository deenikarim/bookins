package handlers

import (
	"context"
	"github.com/deenikarim/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

//postData hold whatever we are posting to a page
type postData struct {
	key   string //input field name
	value string //Values for the input fields
}

//theTests variable is meant for the actual testing(we now have defined our struct but have no values)
//it is a slice of struct because we are going to have more than one test we want to run in out table testing
var theTests = []struct {
	name   string //will be whatever I want to call the individual tests
	url    string // the path which is matched by routes
	method string //whether a GET request or a POST request
	//params             []postData //will be the things that are being posted "a form can have multiple inputs so slice"
	expectedStatusCode int //whether the test has passed or not
}{
	//POPULATE THE STRUCT WITH VALUES
	//because it is a slice of struct, in here each entry have its own curly braces
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"rs", "/reservation-summary", "GET", http.StatusOK},
	/*/for posting handlers
	//for post PostSearchAvailability (actual test data entry)
	{"post-search-avail", "/search-availability", "POST", []postData{
		{"start", "2022-12-6"},
		{"end", "2022-12-12"},
	}, http.StatusOK},

	//for check-availability-json (actual test data entry)
	{"post-check-avail-json", "/search-availability-json", "POST", []postData{
		{"start", "2022-12-6"},
		{"end", "2022-12-12"},
	}, http.StatusOK},

	//for Make reservation (actual test data entry)
	{"make-reservation-post", "/make-reservation", "POST", []postData{
		{"first-name", "John"},
		{"last_name", "kareem"},
		{"email", "john@kareem.com"},
		{"phone", "555-555-5555"},
	}, http.StatusOK}, */
}

//TestHandlers is for testing all the handlers
func TestHandlers(t *testing.T) {
	//get the routes (getRoute() func)
	routes := getRoute()
	//create a new test server
	ts := httptest.NewTLSServer(routes)
	//Now once you create a new test server, it's going to fire up for the life of the test, it's going to listen on
	//port and also open certain things up (so you have to close it when you are done with it)
	defer ts.Close() //anything after defer doesn't get executed until after the current function is finished

	//running the individual test(current table test)
	for _, e := range theTests {
		//what we want to test here
		//we need to make a request as a client as though we were on web browser accessing web page
		resp, err := ts.Client().Get(ts.URL + e.url) //append our url to the test server url
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		//checking our response by looking at the status code we got back from test server and compare it with the expected status code
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

		/*/two different tests we are going to run (GET request or POST request)
		//GET request
		if e.method == "GET" {
			//what we want to test here
			//we need to make a request as a client as though we were on web browser accessing web page
			resp, err := ts.Client().Get(ts.URL + e.url) //append our url to the test server url
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			//checking our response by looking at the status code we got back from test server and compare it with the expected status code
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}

		/*else {
			//holds information as a post request for the variable
			values := url.Values{}

			//And as we through that we will populate values into that values variable
			for _, p := range e.params {
				values.Add(p.key, p.value) //populated values variable with everything necessary to make a post(our data entry)
			}
			//we need to make a request as a client as though we were on web browser accessing web page
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			//checking our response by looking at the status code we got back from test server and compare it with the expected status code
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}*/

	}

}

//TestRepository_Reservation writes test for the Reservation handler
func TestRepository_Reservation(t *testing.T) {
	// ALook up of the reservation handler function, it is expecting to pull a models.Reservation
	// out of the session and store it in a variable called res, so we need to have a models.Reservation to
	// put it in the session

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		}, //now we have reservation variable so need to put it in the session
	}

	//TODO: 1: to put something in a session,  we need a request
	//let's get a request
	req, _ := http.NewRequest("GET", "make-reservation", nil)

	//now we need to get our context into that Request (creating the context)
	ctx := getCtx(req) //so now we have a context that we add to the request

	//add the context to our request
	req = req.WithContext(ctx) //have a request that knows about the X-Session

	//Next: rr for request recorder (have a response recorder)
	// NewRecorder() simulates what we get from the request-response cycle, when someone fires up
	//a web browser, hits our website, gets to a handler then passes it a request, gets a response writer
	//that writes response to the web browser::: NewRecorder() fakes this entire process
	rr := httptest.NewRecorder()

	//put our reservation into the session
	session.Put(ctx, "reservation", reservation)

	//Next: is to create the reservation function which is handler so cannot call it directly
	//TO DO: take the reservation function handler and turn it into handlers so that we call
	handler := http.HandlerFunc(Repo.Reservation)

	//calling that handler
	handler.ServeHTTP(rr, req)

	//At this point, lets determine if our test passed or not
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//summary: HAVE THE MEANS OF PUTTING SOMETHING IN THE SESSION AND PULLING OUT OF THE SESSION

	//TODO: 2: test cases where there is no reservation in the session in the same handler
	//todo: so resetting everything from the above code and ignore reservation variable
	//reset our request by reinitializing it
	req, _ = http.NewRequest("GET", "make-reservation", nil)

	//still need to get the context with the session header because the reason is if we dont do that, then we can not
	//even test the situation where we cannot find a value in the session because there is no session to begin with
	ctx = getCtx(req)

	//the above code gives us a context that we can then put back into the request so now we have session but that session
	//does not have the reservation variable in it because we are not going to add it
	req = req.WithContext(ctx)

	//create a new response recorder
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	//At this point, lets determine if our test passed or not
	//StatusTemporaryRedirect because that is what we expect to find in there
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//TODO: 3: test with non-existing room
	//reset our request by reinitializing it again
	req, _ = http.NewRequest("GET", "make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	//modify the reservation model to make sure that the roomID is greater than two, which implies that there is
	// no such room
	reservation.RoomID = 100
	//put back our session variable reservation into the session so that we can pass the first case
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	//At this point, lets determine if our test passed or not
	//StatusTemporaryRedirect because that is what we expect to find in there
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

//TestRepository_PostReservation writes test for the PostReservation handler
func TestRepository_PostReservation(t *testing.T) {
	//TODO: what we need to do is generate post request that supplies all body that will pass the form data

	//todo: CASE-1: TESTING FOR EVERYTHING SITUATION IS CORRECT THAT IS WE DID EVERYTHING RIGHT

	// let's build a post request to know how post request work
	//inputting valid data
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1")

	//Building up a request which has post data as a body
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	//creates a new context
	//still need to get the context with the session header because the reason is if we don't do that, then we can not
	//even test the situation where we cannot find a value in the session because there is no session to begin with
	ctx := getCtx(req) //have a context that knows about the session

	//add the context to the request
	req = req.WithContext(ctx)

	//may not be required but is good practice to set the header of the request to tell the web server about the
	// kind of request coming its way.
	// TIP: tell the web server, the request you are about to get is form post
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//this satisfies the requirement for a responseWriter
	rr := httptest.NewRecorder()

	//now we need our handler
	handler := http.HandlerFunc(Repo.PostReservation)

	//calling the handler
	handler.ServeHTTP(rr, req)

	//expected return
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//todo: CASE -2: if we cannot parse the parseFrom(),
	// test for missing post body to the request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// TODO: CASE-3 : test for invalid start date
	postedData = url.Values{}
	postedData.Add("start_date", "invalid") //help fail because of the invalid start date entry
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.ca")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// TODO: CASE-4 : test for invalid end date
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "invalid") //help fail because of the invalid end date entry
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.ca")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// TODO: CASE-5: test for invalid room id
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.ca")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "invalid") //help fail because of the invalid roomID entry thus can not be converted to an int

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// TODO: CASE-6: test for invalid data
	// 		making sure that our form pass validation so tests can situations where it fails validation
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "J") //help fail test because first_name required at least three characters long to pass validation
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.ca")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// TODO: CASE -7: test for failure to insert reservation into database
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.ca")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// TODO: CASE-8: test for failure to insert restriction into database
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.ca")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting room restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

//TestRepository_ReservationSummary writes test for reservation summary handler
func TestRepository_ReservationSummary(t *testing.T) {
	// TODO: ALook up of the reservation handler function, it is expecting to pull a models.Reservation
	// 		out of the session and store it in a variable called res, so we need to have a models.Reservation to
	// 		put it in the session
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		}, //now we have reservation variable so need to put it in the session
	}

	//TODO: 1: In case everything works right (first case -- reservation in session)
	//let's get a request
	req, _ := http.NewRequest("GET", "reservation-summary", nil)

	//now we need to get our context into that Request (creating the context)
	ctx := getCtx(req) //so now we have a context that we add to the request

	//add the context to our request
	req = req.WithContext(ctx) //have a request that knows about the X-Session

	//Next: rr for request recorder (have a response recorder)
	// NewRecorder() simulates what we get from the request-response cycle, when someone fires up
	//a web browser, hits our website, gets to a handler then passes it a request, gets a response writer
	//that writes response to the web browser::: NewRecorder() fakes this entire process
	rr := httptest.NewRecorder()

	//put our reservation into the session
	session.Put(ctx, "reservation", reservation)

	//Next: is to create the reservation function which is handler so cannot call it directly
	//TO DO: take the reservation function handler and turn it into handlers so that we call
	handler := http.HandlerFunc(Repo.ReservationSummary)

	//calling that handler
	handler.ServeHTTP(rr, req)

	//At this point, lets determine if our test passed or not
	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//TODO: 2: test cases where there is no reservation in the session in the same handler
	//todo: so resetting everything from the above code and ignore reservation variable
	//reset our request by reinitializing it
	req, _ = http.NewRequest("GET", "reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

//TestRepository_ChooseRoom writes test for the ChooseRoom handler
func TestRepository_ChooseRoom(t *testing.T) {
	/*****************************************
	//todo: first case -- reservation in session
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the RequestURI on the request so that we can grab the ID from the URL
	// In your test for ChooseRoom, you will want to set the URL on your request as follows:
	// req.RequestURI = "/choose-room/1"
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	//expected return for the whole function
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// todo: second case(2) -- reservation not in session
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// todo: third case(3) -- missing url parameter, or malformed parameter
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

//TestRepository_BookRoom writes test for the BookNow handler
func TestRepository_BookRoom(t *testing.T) {
	/*****************************************
	// todo: first case -- database works
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	//create a new request with our URL values grabbed ---id, s, and e
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// todo: second case -- database failed to pull the rooms name and other information
	*****************************************/
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

//TestRepository_PostSearchAvailability writes test for the post search availability handler
func TestRepository_PostSearchAvailability(t *testing.T) {

	// TODO: first case -- rooms are not available

	// create our request body
	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-03")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.PostSearchAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostSearchavailability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// TODO: second case -- rooms are available by the database function that does search for availability for rooms

	// this time, we specify a start date before 2040-01-01, which will give us
	// a non-empty slice, indicating that rooms are available
	postedData = url.Values{}
	postedData.Add("start", "2040-01-01")
	postedData.Add("end", "2040-01-02")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostSearchAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("PostSearchAvailability when rooms are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// todo: third case -- empty post data body

	// create our request with a nil body, so parsing form fails
	req, _ = http.NewRequest("POST", "/search-availability", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostSearchAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// todo: fourth case -- start date in wrong format

	// this time, we specify a start date in the wrong format
	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2040-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostSearchAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// todo: fifth case (5) -- end date in wrong format

	// this time, we specify a start date in the wrong format
	postedData = url.Values{}
	postedData.Add("start", "2040-01-01")
	postedData.Add("end", "invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostSearchAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// todo: sixth case -- database query fails

	// this time, we specify a start date of 2060-01-01, which will cause
	// our testdb repo to return an error
	postedData = url.Values{}
	postedData.Add("start", "2060-01-01")
	postedData.Add("end", "2060-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostSearchAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

/*
func TestRepository_AvailabilityJSON(t *testing.T) {

	// first case -- rooms are not available

	// create our request body
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	// this time we want to parse JSON and get the expected response
	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date > 2049-12-31, we expect no availability
	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}


	// second case -- rooms not available

	// create our request body
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if !j.OK {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}


	// third case -- no request body

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Internal server error" {
		t.Error("Got availability when request body was empty")
	}



	// fourth case -- database error

	// create our request body /reqBody = "start=2060-01-01"

	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Error querying database" {
		t.Error("Got availability when simulating database error")
	}
}
*/

// getCtx creates a new context
func getCtx(req *http.Request) context.Context {
	//create a context
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session")) //get header from the request
	if err != nil {
		log.Println(err)
	}
	return ctx
}
