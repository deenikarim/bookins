package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

//TestForm_Validate write test for the Validate() function
func TestForm_Validate(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST", "/whatever", nil)

	//create an empty form
	form := New(r.PostForm)

	//then just call form.validate() and store the result in the variable "isValidate"
	isValidate := form.Validate()

	//so since there is nothing in the form to call it, to fail validation
	//we just check to see if isValidate is false and then throw an error
	if !isValidate {
		t.Error("got invalid when should have been valid")
	}
}

//TestForm_Required write test for the Required() function
func TestForm_Required(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST", "/whatever", nil)

	//create an empty form
	form := New(r.PostForm)

	//for this form it required the fields of A B C
	form.Required("a", "b", "c")
	//if this form is validate() then it should fail the test because it required inputs for some given field
	//but the values are missing
	if form.Validate() {
		t.Error("form shows valid when required fields are missing") //failure signal
	}

	//created variable called postedData to contains our form values and is of type url.Values
	//and we added values to it
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	//calling the request method again and this time overriding with new value(a new request)
	r, _ = http.NewRequest("POST", "/whatever", nil)

	//set the PostForm field on that request object to postedData
	r.PostForm = postedData
	//so this time our PostForm now has values defined in postedData variable(reinitializing the form)
	form = New(r.PostForm)
	//Second test:for this form it required the fields of A B C
	form.Required("a", "b", "c")
	//if this form is not validate() then it should fail the test
	if !form.Validate() {
		t.Error("shows does not have required fields when it does") //failure signal
	}

}

//TestForm_Has writes test for the Has() function
func TestForm_Has(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST", "/whatever", nil)

	//create an empty form
	form := New(r.PostForm)

	//now need to call the Has() function and get tested
	//code happen: when we call the Has variable, we are calling it with a request that we made from "r"
	//and then we are posting or adding to the form the post form values from that request which currently empty
	//first test:
	has := form.Has("whatever")

	//above: so when we call the has variable, we are passing nothing in the PostForm so the field should
	//not exist therefore "has" should return false because it doesn't have it

	if has {
		//if it returns true then it is lying because nothing is passed
		t.Error("form shows has field when it does not")
	}

	//NOW; lets create another request and check the situation where we know we have posted values in there
	//created variable called postedData to contains our form values and is of type url.Values
	//and we added values to it
	postedData := url.Values{}
	postedData.Add("a", "a")
	//so this time our PostForm now has values defined in postedData variable(reinitializing the form)
	form = New(postedData)

	//Second test:
	//check this request for existence of field called "a" in the post request
	has = form.Has("a")

	//now should be able to says if it doesn't have it(field called "a")
	if !has {
		//if it returns false
		t.Error("form does not have field when it should")
	}
}

//TestForm_MinLength writes test for the MinLength() function
func TestForm_MinLength(t *testing.T) {
	//create a new request
	//r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}

	//create an empty form
	form := New(postedData)

	//first test: make sure that MinLength doesn't work for a non-existent field
	//no field exists in the PostForm right now
	form.MinLength("x", 10)
	//we know form.validate should return false because we are checking for the Minlength of a non-existent field
	if form.Validate() {
		//this if condition should be set to false(so fail the test)
		t.Error("form shows minimum length for non-existent field")
	}

	//using the Get() function
	//find somewhere where we can call Get() function and get an error so if the code above fails then
	//we should get an error
	isError := form.Errors.Get("x")
	if isError == "" {
		//return an error because didn't work
		t.Error("should have an error but got none")
	}

	//NOW; lets create another request and check the situation where we know we have posted values in there
	//created variable called postedData to contains our form values and is of type url.Values
	//and we added values to it
	postedData = url.Values{}
	postedData.Add("some_field", "some_value")

	//so this time our PostForm now has values defined in postedData variable(reinitializing the form)
	//this creates a new form and passed it the information above for the form values
	form = New(postedData)

	//Second test:
	//now field exists in the PostForm right now
	form.MinLength("some_field", 100)
	//we know form.validate should return false because the value is less than 100 characters long
	if form.Validate() {
		//this if condition should be set to false(so fail the test)
		t.Error("shows minimum length of 100 met when data is shorter")
	}

	//third test:
	//reinitializing our postedData to empty and adding values to it
	postedData = url.Values{}
	postedData.Add("another_field", "123abc")

	//now let's reinitialize our form
	form = New(postedData)
	form.MinLength("another_field", 1)

	//if it doesn't pass validation then fail the test, but we know this test should passed
	if !form.Validate() {
		//if it returns false
		t.Error("shows minimum length of 1 met when it is")
	}

	//using the Get() function
	//find somewhere where we can call Get() function and get an error so if the code above fails then
	//we should get an error
	isError = form.Errors.Get("another_field")
	//if error is not equal to, nothing meaning there is some value in there
	if isError != "" {
		t.Error("should not have an error but got one")
	}

}

//TestForm_IsEmail writes test for the IsEmail() function
func TestForm_IsEmail(t *testing.T) {
	//create a new request
	//r := httptest.NewRequest("POST", "/whatever", nil) //we don't need it(useless)

	postedData := url.Values{}

	//create an empty form
	form := New(postedData)

	//passed adding non-existent field
	form.IsEmail("x")
	//if the condition is true or valid then fail the test ( it can't be valid because we passed it no value)
	if form.Validate() {
		t.Error("form shows valid email for a non-existent field")
	}

	//VALID EMAIL ADDRESS
	//created variable called postedData to contains our form values and is of type url.Values
	//and we added values to it
	postedData = url.Values{}
	postedData.Add("email", "me@gmail.com")

	//so this time our PostForm now has values defined in postedData variable(reinitializing the form)
	//this creates a new form and passed it the information above for the form values
	form = New(postedData)

	form.IsEmail("email")

	//now we should be getting a valid email address because we passed it a field name email with a value in
	// the form of an email address
	if !form.Validate() {
		//if it's not valid email address it should fail but since we have a valid email address it should
		//pass the test
		t.Error("got invalid email address when should not be getting that")
	}

	//INVALID EMAIL ADDRESS
	postedData = url.Values{}
	postedData.Add("email", "me@")

	//so this time our PostForm now has values defined in postedData variable(reinitializing the form)
	//this creates a new form and passed it the information above for the form values
	form = New(postedData)

	form.IsEmail("email")

	//if the form is valid this time then fail the test(result is an invalid email address)
	if form.Validate() {
		t.Error("got valid for invalid email address when should not be getting that")
	}
}
