package main

import (
	"errors"
	"fmt"
	"net/http"
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
	if err := app.groupRepository.GetByID(uint(groupId), &group, "Users"); errors.Is(err, gorm.ErrRecordNotFound) {
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
	form.IsDatabaseID("select-owner")
	form.MaxLength("amount", 5)
	form.IsFloat64("amount")
	form.MaxLength("select-participants", 100)

	if !form.Valid() {
		app.render(w, r, "createExpense.page.template.html", &templateData{
			Form:    form,
			GroupID: groupId,
		})
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
