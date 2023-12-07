package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors validationErrors
}

func New(data url.Values) *Form {
	return &Form{data, validationErrors(map[string][]string{})}
}

func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

func (form *Form) MaxLength(field string, maxCount int) {
	value := form.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > maxCount {
		form.Errors.Add(field, fmt.Sprintf("This field is too long (%d max)", maxCount))
	}
}

func (form *Form) PermittedValues(field string, opts ...string) {
	value := form.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	form.Errors.Add(field, "This value is invalid")
}

func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}
