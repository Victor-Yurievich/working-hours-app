package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserToBlock struct {
	Username string `json:"username"`
}

type LastRun struct {
	Time string `json:"time"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Println("Calling " + r.URL.Path)
		fn(w, r)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, settings *Settings) {
	err := templates.ExecuteTemplate(w, tmpl+".html", settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseLoginRequestForm(r *http.Request) (username string, password string) {
	username = r.FormValue("username")
	password = r.FormValue("password")
	return username, password
}

func parseSettingsRequestForm(r *http.Request, w http.ResponseWriter) (from int, to int) {
	r.ParseForm()
	for key, value := range r.Form {
		if key == "from" {
			i, err := strconv.Atoi(value[0])
			handleDecodingError(err, w)
			from = i
		}
		if key == "to" {
			i, err := strconv.Atoi(value[0])
			handleDecodingError(err, w)
			to = i
		}
	}
	return from, to
}

//func populateKey(key int, value string, w http.ResponseWriter) int {
//	i, err := strconv.Atoi(value)
//	handleDecodingError(err, w)
//	key = i
//	return key
//}

func getLoginAttempt(username string, r *http.Request) Login {
	user := getUserByUsername(username)
	loginHour := time.Now().Hour()

	login := Login{
		LoginId:        strconv.Itoa(rand.Intn(1000000)),
		Username:       user.Username,
		UserEmail:      user.Email,
		Ip:             getIP(r),
		UserRole:       user.Role,
		LoginDate:      time.Now().Format("2006-01-02T15:04:05.000Z"),
		LoginHour:      loginHour,
		ValidLoginHour: validateLoginHour(loginHour),
	}
	return login
}

func getUserByUsername(username string) User {
	var retutnUser User
	for _, user := range users {
		if username == user.Username {
			retutnUser = user
		}
	}
	return retutnUser
}

func getIP(r *http.Request) string {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}

func setAuthCookie(w http.ResponseWriter) http.ResponseWriter {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "token", Value: "authenticated", Expires: expiration, HttpOnly: true, Path: "/dashboard"}
	http.SetCookie(w, &cookie)
	return w
}

func deleteAuthCookie(w http.ResponseWriter) http.ResponseWriter {
	expiration := time.Unix(0, 0)
	cookie := http.Cookie{Name: "token", Value: "", Expires: expiration, HttpOnly: true, Path: "/dashboard"}
	http.SetCookie(w, &cookie)
	return w
}

func createPongJsonResponce() []byte {
	response := make(map[string]bool)
	response["is_alive"] = true
	return createJson(response)
}

func createIncidentsResponse(lastRun LastRun) []byte {
	outOfWorkingHoursLogins := retrieveOutOfWorkingHoursLogins(logins, lastRun)
	return createJson(outOfWorkingHoursLogins)
}

func retrieveOutOfWorkingHoursLogins(logins []Login, lastRun LastRun) []Login {
	var loginsToRetrieve []Login
	for _, login := range logins {
		if validateLoginForIncidentResponse(login, lastRun) {
			loginsToRetrieve = append(loginsToRetrieve, login)
		}
	}
	return loginsToRetrieve
}

func createJson(data interface{}) []byte { // Ask Lior about using interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	return jsonData
}

func decodeUserRequestBody(r *http.Request, structObject *UserToBlock) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&structObject)
	return err
}

func decodeUserIncidentsBody(r *http.Request, structObject *LastRun) error {
	decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields()
	err := decoder.Decode(&structObject)
	return err
}

func handleDecodingError(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		return
	}
	errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	response := make(map[string]string)
	response["message"] = message
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}
