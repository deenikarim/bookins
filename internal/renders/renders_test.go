package renders

import (
	"github.com/deenikarim/bookings/internal/models"
	"net/http"
	"testing"
)

//TestAddDefaultData write test for AddDefaultData
func TestAddDefaultData(t *testing.T) {
	//call the first argument
	var td models.TemplateData

	//call getSession() function here
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	//need something to test so lets put something into the session for example our flash
	session.Put(r.Context(), "flash", "123")

	//Now we call the function which is (AddDefaultData)
	result := AddDefaultData(&td, r)

	//doing something with the result variable(checking)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}

}

//TestRenderTemplate write test for our RenderTemplate function
func TestRenderTemplate(t *testing.T) {
	//correct path to templates for testing within render package
	pathToTemplate = "./../../templates"

	//need to have a template cache before we can render anything so lets get our template cache
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	//after creating the template cache we need to put that into our variable "app"
	app.TemplateCache = tc

	//get a request
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	//using our myWriter instead ResponseWriter function
	var ww myWriter

	//call the RenderTemplate
	err = RenderTemplate(&ww, r, "majors.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error rendering template to the browser")
	}

	//checking for something that is not in the template cache(non-existent)
	err = RenderTemplate(&ww, r, "non-existent.page.html", &models.TemplateData{})
	//if it successfully renders a non-existent template, well, it should not be doing that
	if err == nil {
		t.Error("rendered template to the browser that doesn't exist")
	}

}

//TestNewTemplates write test for our NewTemplates function
func TestNewTemplates(t *testing.T) {
	//code: this doesn't do anything so this should be sufficient to go through testing process
	NewTemplates(app)
}

//TestCreateTemplateCache write test for our CreateTemplateCache function
func TestCreateTemplateCache(t *testing.T) {

	//correct location of our templates folder
	pathToTemplate = "./../../templates"

	//getting the template cache
	//we don't need to the tc variable we just need it to run our test
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

//getSession get session data
//solution to the error is to change our test so that it has a request that has session data
func getSession() (*http.Request, error) {
	//second argument which is a request(this creates a New request for us)
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	//we have a request object, but we can't just return it, or we will right back to no session available
	//so create a context, and we use that anytime we write to or read from the session
	ctx := r.Context()

	//so now we need to put the session data into that(put session data into our context)
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session")) //this part make (r.Header.Get("X-Session")) it active session

	//putting our context back into our request
	r = r.WithContext(ctx)

	//returning our request and our error
	return r, err
}
