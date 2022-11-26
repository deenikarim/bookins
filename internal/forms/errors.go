package forms

//

type errors map[string][]string

//Add adds an error message for a given form field
func (e errors) Add(field, message string) {
	//what I am putting in here... is an error for a given field like firstName and something kind of
	//message like firstName must least be three characters
	e[field] = append(e[field], message)
}

//Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field] //assign the values to the whatever we finds in our map[string] with the index of the field
	if len(es) == 0 {
		return ""
	}
	return es[0]
	//code explanation: given the means of serve whether the field given field has a value error
}
