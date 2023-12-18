package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	form.Required("name", "participants")
	form.MaxLength("name", 20)
	form.MaxLength("participants", 100)

	if !form.Valid() {
		app.render(w, r, "createGroup.page.template.html", &templateData{Form: form})
		return
	}

	name := form.Get("name")
	participantsList := strings.Split(form.Get("participants"), ",")
	date, _ := time.Parse("2006-01-02", form.Get("date"))

	var participants []models.User
	dbErrorFindingUsers := app.userRepository.FindByNames(participantsList, &participants).Error
	if errors.Is(dbErrorFindingUsers, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if dbErrorFindingUsers != nil {
		app.serverError(w, dbErrorFindingUsers)
		return
	}

	groupOwner := app.authenticatedUser(r)

	var participantsPtrs []*models.User
	for _, participant := range participants {
		participantsPtrs = append(participantsPtrs, &participant)
	}

	group, errorConstructingGroup := models.NewGroup(name, groupOwner, participantsPtrs, date)

	if errorConstructingGroup != nil {
		app.errorLog.Fatal(errorConstructingGroup)
	}

	insertGroupResponse := app.groupRepository.Create(group)

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

func (app *Application) createExpenseForm(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.
		Query().Get(":id"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	var group models.Group
	dbError := app.groupRepository.GetByID(uint(groupId), &group, "Users").Error
	if errors.Is(dbError, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if dbError != nil {
		app.serverError(w, dbError)
		return
	}

	app.render(w, r, "createExpense.page.template.html", &templateData{
		Group:   &group,
		Form:    forms.New(nil),
		GroupID: groupId,
	})
}

func (app *Application) createExpense(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.
		Query().Get(":id"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	errorParsingForm := r.ParseForm()
	if errorParsingForm != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("description", "amount", "select-participants", "select-owner")
	form.MaxLength("description", 20)
	form.MaxLength("select-owner", 20)
	form.MaxLength("amount", 5)
	form.IsFloat64("amount")
	form.MaxLength("select-participants", 100)

	if !form.Valid() {
		app.render(w, r, "createExpense.page.template.html", &templateData{
			Form:    form,
			GroupID: groupId,
		})
	}

	var expenseOwner models.User
	getOwnerByNameError := app.userRepository.GetByName(form.Get("select-owner"), &expenseOwner).Error
	if getOwnerByNameError != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	float64Amount, _ := strconv.ParseFloat(form.Get("amount"), 64)

	participantsSlice := form.Values["select-participants"]
	var participants []models.User
	var participantsPtrs []*models.User
	errorGettingParticipants := app.userRepository.FindByNames(participantsSlice, &participants).Error
	if errors.Is(errorGettingParticipants, gorm.ErrRecordNotFound) {
		app.clientError(w, http.StatusBadRequest)
		return
	} else if errorGettingParticipants != nil {
		app.serverError(w, errorGettingParticipants)
		return
	}

	for _, participant := range participants {
		tempParticipant := participant
		participantsPtrs = append(participantsPtrs, &tempParticipant)
	}

	expense, errorConstructingExpense := models.NewExpense(form.Get("description"), float64Amount, uint(groupId), expenseOwner.ID, participantsPtrs)
	if errorConstructingExpense != nil {
		app.errorLog.Fatal(errorConstructingExpense)
	}

	insertExpenseError := app.expenseRepository.Create(expense).Error
	if insertExpenseError != nil {
		app.serverError(w, insertExpenseError)
		return
	}

	app.session.Put(r, "flash", "Expense created!")

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", groupId), http.StatusSeeOther)
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
	dbError := app.groupRepository.GetByID(uint(id), group, "Users", "Expenses.Owner", "Expenses.Participants", "Transactions").Error

	if errors.Is(dbError, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if dbError != nil {
		app.serverError(w, dbError)
		return
	}

	app.render(w, r, "showGroup.page.template.html", &templateData{Group: group})
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
