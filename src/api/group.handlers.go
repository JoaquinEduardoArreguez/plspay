package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JoaquinEduardoArreguez/plspay/package/forms"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"gorm.io/gorm"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

func (app *Application) createGroupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createGroup.page.template.html", &templateData{
		Form: forms.New(nil),
	})
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
	participantsNames := strings.Split(form.Get("participants"), ",")
	date, _ := time.Parse("2006-01-02", form.Get("date"))
	groupOwner := app.authenticatedUser(r)

	group, errorCreatingGroup := app.groupService.CreateGroup(name, groupOwner, participantsNames, date)
	if errorCreatingGroup != nil {
		app.errorLog.Fatal()
	}

	app.session.Put(r, "flash", "Group created!")

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", group.ID), http.StatusSeeOther)
}

func (app *Application) showGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	group := &models.Group{}
	if err := app.groupRepository.GetByID(uint(id), group, "Users", "Expenses.Owner", "Expenses.Participants", "Transactions"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "showGroup.page.template.html", &templateData{Group: group})
}

func (app *Application) showGroups(w http.ResponseWriter, r *http.Request) {
	var userGroupsDtos []*models.GroupDTO

	user := app.authenticatedUser(r)
	if err := app.userRepository.GetByID(user.ID, user, "Groups"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
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
