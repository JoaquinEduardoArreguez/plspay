package forms

type validationErrors map[string][]string

func (errors validationErrors) Add(field, message string) {
	errors[field] = append(errors[field], message)
}

func (errors validationErrors) Get(field string) []string {
	return errors[field]
}
