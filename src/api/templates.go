package main

import (
	"html/template"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/JoaquinEduardoArreguez/plspay/package/forms"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
)

type templateData struct {
	Form              *forms.Form
	Groups            []*models.Group
	Flash             string
	AuthenticatedUser *models.User
	CsrfToken         string
	Group             *models.Group
	GroupID           int
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.template.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.template.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.template.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006")
}

func participantsSelectLen(value interface{}) string {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Slice {
		lenght := len(value.([]*models.User))
		if lenght > 4 {
			lenght = 4
		}
		return strconv.Itoa(lenght)
	}
	return ""
}

var functions = template.FuncMap{
	"humanDate":             humanDate,
	"participantsSelectLen": participantsSelectLen,
}
