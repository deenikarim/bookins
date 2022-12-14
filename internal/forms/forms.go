package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

//(working server side validation)

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

//Required checks for required fields and add an error if there are values missing in there
func (f *Form) Required(fields ...string) {
	//range through that field variable which might at least have one thing in it
	for _, field := range fields {
		//what do I want to do here?
		//let's go get the value
		value := f.Get(field) //go to form object and call the method func called GET and get the field name "firstName"

		//if that form has no value then add an error and a message
		if strings.TrimSpace(value) == "" {
			//add an error
			f.Errors.Add(field, "This field can't be blank!")
			//NB: if this condition is never met, then adding of errors will not be fulfilled
		}

	}
}

//Has checks if form fields is post and not empty(check if a field exists)
//defined certain Checks to see if the form data we received is valid or not
func (f *Form) Has(field string) bool {
	//checking to see if a field is fulfilled or exists
	//code:this check a form like a post request for the value of first_name as does it have an entry for
	//first_name, if it has one, and it is an empty string then return false otherwise return true

	//MODIFIED: To check whether the checkbox requirement  is fulfilled
	//x := r.Form.Get(field)
	x := f.Get(field)
	if x == "" {
		//adding errors to input fields if input field is empty and a message
		//f.Errors.Add(field, "This field can't be blank!")
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

//MinLength checks for a field or strings minimum length
func (f *Form) MinLength(field string, length int) bool {
	//getting the value of a field or strings
	//x := r.Form.Get(field)
	x := f.Get(field)
	//checking to see if it is long enough
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true

}

//IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
