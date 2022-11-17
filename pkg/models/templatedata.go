package models

//PART-6: DATA SHARING WITH THE TEMPLATES

//TemplateData hold data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRTToken string
	Flash     string
	Error     string
}
