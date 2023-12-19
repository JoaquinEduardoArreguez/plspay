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
