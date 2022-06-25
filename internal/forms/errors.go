package forms

type errors map[string][]string

// Add adds error message for given field
func (e errors) AddError(field, message string) {
	e[field] = append(e[field], message)
}

//GetError returns the first error message
func (e errors) GetError(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
