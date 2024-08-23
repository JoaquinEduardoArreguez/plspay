package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/JoaquinEduardoArreguez/plspay/package/forms"
	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"gorm.io/gorm"
)

func (app *Application) createExpenseForm(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.
		Query().Get(":id"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	var group models.Group
	if err := app.groupService.GetGroupById(&group, groupId, "Users"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "createExpense.page.template.html", &templateData{
		Group:   &group,
		Form:    forms.New(nil),
		GroupID: groupId,
	})
}

func (app *Application) createExpense(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get(":id"))
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
	form.IsDatabaseID("select-owner")
	form.MaxLength("amount", 20)
	form.IsFloat64("amount")
	form.MaxLength("select-participants", 100)

	var group models.Group
	if err := app.groupService.GetGroupById(&group, groupId, "Users"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "createExpense.page.template.html", &templateData{
			Form:    form,
			GroupID: groupId,
			Group:   &group,
		})
		return
	}

	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	ownerId, _ := strconv.Atoi(form.Get("select-owner"))

	var participantsIds []uint
	for _, participantId := range form.Values["select-participants"] {
		id, _ := strconv.Atoi(participantId)
		participantsIds = append(participantsIds, uint(id))
	}

	_, errorCreatingExpense := app.expenseService.CreateExpense(
		form.Get("description"),
		amount,
		uint(groupId),
		uint(ownerId),
		participantsIds,
	)

	if errorCreatingExpense != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Expense created!")

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", groupId), http.StatusSeeOther)
}

func (app *Application) deleteExpense(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get(":groupId"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	expenseId, err := strconv.Atoi(r.URL.Query().Get(":expenseId"))
	if err != nil || expenseId < 1 {
		http.NotFound(w, r)
		return
	}

	if err := app.expenseService.DeleteExpense(uint(groupId), uint(expenseId)); err != nil {
		app.session.Put(r, "flash", "Error deleting expense")
	}

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", groupId), http.StatusSeeOther)
}

func (app *Application) updateExpenseForm(w http.ResponseWriter, r *http.Request) {
	expenseId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || expenseId < 1 {
		http.NotFound(w, r)
		return
	}

	var expense models.Expense
	if err := app.expenseService.GetExpenseById(&expense, expenseId, "Participants"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	var group models.Group
	if err := app.groupService.GetGroupById(&group, int(expense.GroupID), "Users"); errors.Is(err, gorm.ErrRecordNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	currentData := url.Values{
		"description": {expense.Description},
		"amount":      {strconv.FormatFloat(expense.Amount, 'f', 0, 64)},
	}

	app.render(w, r, "updateExpense.page.template.html", &templateData{
		Group:   &group,
		Form:    forms.New(currentData),
		GroupID: int(expense.GroupID),
		Expense: &expense,
	})
}

func (app *Application) updateExpense(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get(":groupId"))
	if err != nil || groupId < 1 {
		http.NotFound(w, r)
		return
	}

	expenseId, err := strconv.Atoi(r.URL.Query().Get(":expenseId"))
	if err != nil || expenseId < 1 {
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
	form.IsDatabaseID("select-owner")
	form.MaxLength("amount", 20)
	form.IsFloat64("amount")
	form.MaxLength("select-participants", 100)

	if !form.Valid() {
		app.render(w, r, "updateExpense.page.template.html", &templateData{
			Form: form,
		})
		return
	}

	amount, _ := strconv.ParseFloat(form.Get("amount"), 64)
	ownerId, _ := strconv.Atoi(form.Get("select-owner"))

	var participantsIds []uint
	for _, participantId := range form.Values["select-participants"] {
		id, _ := strconv.Atoi(participantId)
		participantsIds = append(participantsIds, uint(id))
	}

	_, errorCreatingExpense := app.expenseService.UpdateExpense(
		uint(expenseId),
		form.Get("description"),
		amount,
		uint(groupId),
		uint(ownerId),
		participantsIds,
	)

	if errorCreatingExpense != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Expense created!")

	http.Redirect(w, r, fmt.Sprintf("/groups/%d", groupId), http.StatusSeeOther)
}
