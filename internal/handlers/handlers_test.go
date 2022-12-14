package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	name               string     //will be whatever I want to call the individual tests
	url                string     // the path which is matched by routes
	method             string     //whether a GET request or a POST request
	params             []postData //will be the things that are being posted "a form can have multiple inputs so slice"
	expectedStatusCode int        //whether the test has passed or not
}{
	//POPULATE THE STRUCT WITH VALUES
	//because it is a slice of struct, in here each entry have its own curly braces
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"rs", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	//for posting handlers
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
		{"address", "Madina-Accra"},
		{"address_two", "UN Street"},
		{"city", "Accra"},
	}, http.StatusOK},
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
		//two different tests we are going to run (GET request or POST request)
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

		} else {
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
		}

	}

}
