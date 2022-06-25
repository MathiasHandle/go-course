package models

import "github.com/mathiashandle/go-course/internal/forms"

// Holds data sent from handlers to templates
type TemplateData struct {
	StringMap      map[string]string
	IntMap         map[string]int
	FloatMap       map[string]float32
	Data           map[string]interface{}
	CSRFToken      string
	FlashMessage   string
	WarningMessage string
	ErrorMessage   string
	Form           *forms.Form
}
