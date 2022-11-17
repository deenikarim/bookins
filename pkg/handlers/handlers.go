package handlers

import (
	"github.com/deenikarim/bookings/pkg/config"
	"github.com/deenikarim/bookings/pkg/models"
	"github.com/deenikarim/bookings/pkg/renders"
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
	renders.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

//About create the about page handler function
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{})

}

//Reservation create the reservation page handler function
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "hello again, says by kareem"

	//calling the renderTemplate function inside the handler function to render the home page to the browser
	renders.RenderTemplate(w, "reservation.page.tmpl", &models.TemplateData{
		//send data to the template
		StringMap: stringMap,
	})
}

//*********************************************** END ********************************************//
