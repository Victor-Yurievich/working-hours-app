package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(login|auth|dashboard|settings|logout|ws|log-user-out|ping|fetch-incidents)$")
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
	validLogin := validateLoginAttempt(username, password)
	if validLogin != true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	login := getLoginAttempt(username, r)
	logins = append(logins, login)
	err := saveLogin(&logins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w = setAuthCookie(w)
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	fmt.Println(cookies)
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
	renderTemplate(w, "dashboard", &settings)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	from, to := parseSettingsRequestForm(r)
	settings.From = from
	settings.To = to
	err := saveSettings(&settings)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	incidentResponse := createIncidentsResponse()
	w.Write(incidentResponse)
}
