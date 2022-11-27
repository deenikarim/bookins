package forms

//errors will hold our errors which is a map of type string with every value being a slice of strings
//because we might have more than one error for a given field in a form (why is slice)
type errors map[string][]string

//Add adds an error message for a given form field
//we have a function to add errors which will append an error and add message and associate it to a
//particular field
func (e errors) Add(field, message string) {
	//what I am putting in here... is an error for a given field like firstName and something kind of
	//message like firstName must least be three characters
	e[field] = append(e[field], message)
}

//Get returns the first error message
// we also have a function to checks and see if there are errors and returns there first error if there is one
func (e errors) Get(field string) string {
	es := e[field] //assign the values to the whatever we find in our map[string] with the index of the field
	if len(es) == 0 {
		return ""
	}
	return es[0]
	//code explanation: given the means of serve whether the field given field has a value error
}
