package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

func (form *Form) MinLength(field string, count int) {
	value := form.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < count {
		form.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d)", count))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (form *Form) IsFloat64(field string) {
	value := form.Get(field)
	if value == "" {
		return
	}

	_, errorParsingFloat := strconv.ParseFloat(value, 64)
	if errorParsingFloat != nil {
		form.Errors.Add(field, "This field must be a number")
	}
}

func (form *Form) IsDatabaseID(field string) {
	value := form.Get(field)
	if value == "" {
		return
	}

	_, errorParsingFloat := strconv.ParseUint(value, 10, 64)
	if errorParsingFloat != nil {
		form.Errors.Add(field, "This field is invalid")
	}
}
