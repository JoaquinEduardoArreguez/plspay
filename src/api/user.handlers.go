package main

import (
	"errors"
	"net/http"

	"github.com/JoaquinEduardoArreguez/plspay/package/forms"
	"github.com/JoaquinEduardoArreguez/plspay/package/models/repositories"
	"gorm.io/gorm"
)

func (app *Application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.template.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) signupUser(w http.ResponseWriter, r *http.Request) {
	parseFormError := r.ParseForm()
	if parseFormError != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRx)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.template.html", &templateData{Form: form})
		return
	}

	_, err := app.userRepository.CreateUser(form.Get("name"), form.Get("email"), form.Get("password"))

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		form.Errors.Add("email", "Email address already in use")
		app.render(w, r, "signup.page.template.html", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Signup successful, please log in.")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *Application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.template.html", &templateData{Form: forms.New(nil)})
}

func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	id, err := app.userRepository.Authenticate(form.Get("email"), form.Get("password"))
	if err == repositories.InvalidCredentialsError {
		form.Errors.Add("generic", "Invalid credentials")
		app.render(w, r, "login.page.template.html", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	app.session.Put(r, "flash", "Logged in successfully!")

	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

func (app *Application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
