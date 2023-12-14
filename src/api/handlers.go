package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JoaquinEduardoArreguez/plspay/package/forms"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/models/repositories"
	"gorm.io/gorm"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

func (app *Application) createGroup(w http.ResponseWriter, r *http.Request) {
	errorParsingForm := r.ParseForm()
	if errorParsingForm != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "date")
	form.MaxLength("name", 20)

	if !form.Valid() {
		app.render(w, r, "createGroup.page.template.html", &templateData{Form: form})
		return
	}

	name := form.Get("name")
	date, _ := time.Parse("2006-01-02", form.Get("date"))

	groupOwner := app.authenticatedUser(r)
	group, errorConstructingGroup := models.NewGroup(groupOwner, name, date)

	if errorConstructingGroup != nil {
		app.errorLog.Fatal(errorConstructingGroup)
	}

	insertGroupResponse := app.groupRepository.Create(group) // error ignored

	if insertGroupResponse.Error != nil {
		app.serverError(w, insertGroupResponse.Error)
	}

	app.session.Put(r, "flash", "Group created!")

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", group.ID), http.StatusSeeOther)
}

func (app *Application) createGroupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createGroup.page.template.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Application) groupsForm(w http.ResponseWriter, r *http.Request) {
	var userGroupsDtos []*models.GroupDTO

	user := app.authenticatedUser(r)
	dbResponse := app.userRepository.GetByID(user.ID, user, "Groups")
	if errors.Is(dbResponse.Error, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if dbResponse.Error != nil {
		app.serverError(w, dbResponse.Error)
		return
	}

	for _, group := range user.Groups {
		dto := group.ToDto()
		userGroupsDtos = append(userGroupsDtos, &dto)
	}

	app.render(w, r, "groups.page.template.html", &templateData{
		GroupDtos: userGroupsDtos,
	})

}

func (app *Application) getGroupById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	group := &models.Group{}
	dbResponse := app.groupRepository.DB.Preload("Users").First(group, id)

	if dbResponse.Error != nil {
		switch dbResponse.Error {
		case gorm.ErrRecordNotFound:
			app.notFound(w)
		default:
			app.serverError(w, dbResponse.Error)
		}
		return
	}

	groupDto := group.ToDto()

	app.render(w, r, "show.page.template.html", &templateData{GroupByIdDto: &groupDto})
}

func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	user, errorConstructingUser := models.NewUser("Adri", "a@asd.com")

	if errorConstructingUser != nil {
		app.errorLog.Fatal(errorConstructingUser)
	}

	errorCreatingUser := app.userRepository.Create(user) // error ignored

	if errorCreatingUser != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
	}
}

func (app *Application) createUserForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create user form"))
}

func (app *Application) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	user := &models.User{}
	dbResponse := app.userRepository.DB.Preload("Groups").First(user, id)

	if dbResponse.Error != nil {
		switch dbResponse.Error {
		case gorm.ErrRecordNotFound:
			app.notFound(w)
		default:
			app.serverError(w, dbResponse.Error)
		}
		return
	}

	// Convert the user to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

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
