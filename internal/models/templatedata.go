package models

import "github.com/deenikarim/bookings/internal/forms"

//PART-6: DATA SHARING WITH THE TEMPLATES

//TemplateData hold data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} //store the reservation data here
	CSRFToken string
	Flash     string
	Error     string
	Warning   string
	Form      *forms.Form //Available to every page whether it has a form or not
}
