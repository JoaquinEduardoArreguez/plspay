package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"gorm.io/gorm"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Get groups to display
	var groups []models.Group
	dbResponse := app.groupRepository.DB.Preload("Users").Find(&groups)

	if dbResponse.Error != nil {
		app.serverError(w, dbResponse.Error)
		return
	}

	// Convert them to DTOs
	var groupDtos []*models.GroupDTO
	for _, group := range groups {
		dto := group.ToDto()
		groupDtos = append(groupDtos, &dto)
	}

	app.render(w, r, "home.page.template.html", &templateData{
		GroupDtos: groupDtos,
	})

}

func (app *Application) createGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// group owner (User)
	groupOwner := &models.User{}
	dbResponse := app.userRepository.GetByID(1, groupOwner)

	if dbResponse.Error != nil {
		switch dbResponse.Error {
		case gorm.ErrRecordNotFound:
			app.notFound(w)
		default:
			app.serverError(w, dbResponse.Error)
		}
		return
	}

	// Create Group
	group, errorConstructingGroup := models.NewGroup(groupOwner, "Asado", "Asado 23 jul")

	if errorConstructingGroup != nil {
		app.errorLog.Fatal(errorConstructingGroup)
	}

	errorCreatingGroup := app.groupRepository.Create(group) // error ignored

	if errorCreatingGroup != nil {
		http.Error(w, "Error creating group", http.StatusInternalServerError)
	}
}

func (app *Application) getGroupById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	user, errorConstructingUser := models.NewUser("Adri", "a@asd.com")

	if errorConstructingUser != nil {
		app.errorLog.Fatal(errorConstructingUser)
	}

	errorCreatingUser := app.userRepository.Create(user) // error ignored

	if errorCreatingUser != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
	}
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
