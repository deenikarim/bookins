package forms

import (
	"net/http"
	"net/url"
)

//Form creates a custom form struct, embed a url.Values // will hold our form values
type Form struct {
	url.Values        //holders some values from a form
	Error      errors //variable created in error.go

}

//New initializes the form struct
func New(data url.Values) *Form {
	return &Form{
		data, // data will correspond to the inputs on our form fields
		errors(map[string][]string{}),
	}
}

//Has checks if form fields is post and not empty
//defined certain Checks to see if the form data we received is valid
func (f *Form) Has(field string, r *http.Request) bool {
	//checking to see if a field is fulfilled
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}
