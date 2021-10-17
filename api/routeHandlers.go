package api

import (
	"errors"
	common_types "github.com/Victor-Yurievich/working-hours-app/common-types"
	"github.com/Victor-Yurievich/working-hours-app/model"
	"github.com/Victor-Yurievich/working-hours-app/validators"
	"github.com/Victor-Yurievich/working-hours-app/websockets"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(login|auth|dashboard|settings|logout|ws|block-user|ping|fetch-incidents)$")
var templates = template.Must(template.ParseFiles("./view/login.html", "./view/dashboard.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w = deleteAuthCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	username, password := parseLoginRequestForm(r)
	validLogin := validators.ValidateLoginAttempt(username, password)
	if validLogin != true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	login := getLoginAttempt(username, r)
	model.Logins = append(model.Logins, login)
	err := model.SaveLogin(&model.Logins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w = setAuthCookie(w)
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if strings.Contains(cookie.Name, "token") {
			dashboardHandler(w)
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusUnauthorized)
}

func dashboardHandler(w http.ResponseWriter) {
	w = setAuthCookie(w)
	renderTemplate(w, "dashboard", &model.TimeSettings)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	from, to := parseSettingsRequestForm(r, w)
	model.TimeSettings.From = from
	model.TimeSettings.To = to
	err := model.SaveSettings(&model.TimeSettings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func pongResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponce := createPongJsonResponce()
	w.Write(jsonResponce)
}

func incidentsResponse(w http.ResponseWriter, r *http.Request) {
	var lastRun common_types.LastRun
	err := decodeUserIncidentsBody(r, &lastRun)
	if err != nil {
		handleDecodingError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	incidentResponse := createIncidentsResponse(lastRun)
	w.Write(incidentResponse)
}

func blockUserHandler(w http.ResponseWriter, r *http.Request) {
	var user UserToBlock
	err := decodeUserRequestBody(r, &user) // Ask Lior how to do it Generic
	if err != nil {
		handleDecodingError(err, w)
		return
	}
	handleUsersUpdate(w, user)
}

func handleUsersUpdate(w http.ResponseWriter, user UserToBlock) {
	userUpdateError := UpdateUsers(user)
	if userUpdateError != nil {
		errorResponse(w, userUpdateError.Error(), http.StatusBadRequest)
		return
	}
	processBlockUserResponse(w, user.Username)
}

func processBlockUserResponse(w http.ResponseWriter, username string) {
	websockets.LogUserOut()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	blockUserResponse := createBlockUserResponse(username)
	w.Write(blockUserResponse)
}

func UpdateUsers(userToBlock UserToBlock) error {
	for i, user := range model.Users {
		if user.Username == userToBlock.Username {
			user.Blocked = true
			model.Users[i] = user
			model.SaveUsers()
			return nil
		}
	}
	return errors.New("User " + userToBlock.Username + " not found")
}
