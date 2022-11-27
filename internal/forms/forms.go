package forms

import (
	"net/http"
	"net/url"
)

//Form creates a custom form struct, embed a url.Values // will hold our form values
//describe the form data that we are going to add errors to if anything exists
type Form struct {
	url.Values        //holders some values from a form
	Errors     errors //variable created in error.go

}

//New initializes the form struct
//allows us to create a new form (a new empty form)
func New(data url.Values) *Form {
	return &Form{
		data, // data will correspond to the inputs on our form fields
		errors(map[string][]string{}),
	}
}

//Has checks if form fields is post and not empty(check if a field exists)
//defined certain Checks to see if the form data we received is valid or not
func (f *Form) Has(field string, r *http.Request) bool {
	//checking to see if a field is fulfilled
	//code:this check a form like a post request for the value of first_name as does it have an entry for
	//first_name, if it has one, and it is an empty string then return false otherwise return true
	x := r.Form.Get(field)
	if x == "" {
		//adding errors to input fields if input field is empty and a message
		f.Errors.Add(field, "This field can't be blank'")
		return false
	}
	return true
}

//Validate return true if there are no errors, otherwise false
func (f *Form) Validate() bool {
	//CODE:it has receiver of *form. so that we can look and our form data(( all it does is, Say if  there are
	//anything errors at all associated with the form object "f" that I am getting from my receiver return false otherwise return true)
	return len(f.Errors) == 0

}
