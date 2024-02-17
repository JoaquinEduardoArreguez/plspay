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
	"github.com/google/uuid"
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
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "participants")
	form.MaxLength("name", 40)
	form.MaxLength("participants", 500)

	if !form.Valid() {
		app.render(w, r, "createGroup.page.template.html", &templateData{Form: form})
		return
	}

	name := form.Get("name")
	date, _ := time.Parse("2006-01-02", form.Get("date"))
	groupOwner := app.authenticatedUser(r)
	participantsEmails := strings.Split(form.Get("participants"), ",")
	var registeredParticipantsEmails []string
	var guestUsers []*models.User

	groupUuid := uuid.New()

	for _, participantEmail := range participantsEmails {
		if strings.HasSuffix(participantEmail, "@guest.com") {
			var guestName string
			arrobaIndex := strings.Index(participantEmail, "@")
			if arrobaIndex != -1 {
				guestName = participantEmail[:arrobaIndex]
			}

			guestUsers = append(guestUsers, &models.User{Name: guestName, Email: groupUuid.String() + participantEmail, IsGuest: true})
		} else {
			registeredParticipantsEmails = append(registeredParticipantsEmails, participantEmail)
		}
	}

	group, errorCreatingGroup := app.groupService.CreateGroup(groupUuid, name, groupOwner, registeredParticipantsEmails, date, guestUsers)
	if errorCreatingGroup != nil {
		app.serverError(w, errorCreatingGroup)
	}

	app.session.Put(r, "flash", "Group created!")

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", group.ID), http.StatusSeeOther)
}

func (app *Application) showGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	var group models.Group

	if err := app.groupService.GetGroupById(
		&group,
		groupId,
		"Users",
		"Expenses.Owner",
		"Expenses.Participants",
		"Transactions.Sender",
		"Transactions.Receiver",
	); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "showGroup.page.template.html", &templateData{Group: &group})
}

func (app *Application) showGroups(w http.ResponseWriter, r *http.Request) {
	user := app.authenticatedUser(r)

	if err := app.userService.GetUserById(user, int(user.ID), "Groups"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "groups.page.template.html", &templateData{Groups: user.Groups})
}

func (app *Application) calculateTransactions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	transactions, errorCalculatingTransactions := app.groupService.CreateTransactions(uint(id))
	if errorCalculatingTransactions != nil {
		app.serverError(w, errorCalculatingTransactions)
		return
	}

	app.session.Put(r, "flash", fmt.Sprintf("%d transactions calculated", len(transactions)))

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", id), http.StatusSeeOther)
}

func (app *Application) deleteGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	if err := app.groupService.DeleteGroup(uint(groupId)); err != nil {
		app.session.Put(r, "flash", "Error deleting group")
	}

	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}
