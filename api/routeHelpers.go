package api

import (
	"encoding/json"
	"errors"
	"fmt"
	common_types "github.com/Victor-Yurievich/working-hours-app/common-types"
	"github.com/Victor-Yurievich/working-hours-app/model"
	"github.com/Victor-Yurievich/working-hours-app/validators"
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

type UserToBlockReturned struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Blocked  bool   `json:"blocked"`
}

//type LastRun struct {
//	Time string `json:"time"`
//}

//type LastRun struct {
//	Time string `json:"time"`
//}

//Ask Lior about common_types interoperability

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

func renderTemplate(w http.ResponseWriter, tmpl string, settings *model.Settings) {
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

//Ask Lior about generic functions

//func populateKey(key int, value string, w http.ResponseWriter) int {
//	i, err := strconv.Atoi(value)
//	handleDecodingError(err, w)
//	key = i
//	return key
//}

func getLoginAttempt(username string, r *http.Request) model.Login {
	user := getUserByUsername(username)
	loginHour := time.Now().Hour()

	login := model.Login{
		LoginId:        strconv.Itoa(rand.Intn(1000000)),
		Username:       user.Username,
		UserEmail:      user.Email,
		Ip:             getIP(r),
		UserRole:       user.Role,
		LoginDate:      time.Now().Format("2006-01-02T15:04:05.000Z"),
		LoginHour:      loginHour,
		ValidLoginHour: validators.ValidateLoginHour(loginHour),
	}
	return login
}

func getUserByUsername(username string) model.User {
	var retutnUser model.User
	for _, user := range model.Users {
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

func createIncidentsResponse(lastRun common_types.LastRun) []byte {
	outOfWorkingHoursLogins := retrieveOutOfWorkingHoursLogins(model.Logins, lastRun)
	return createJson(outOfWorkingHoursLogins)
}

func createBlockUserResponse(username string) []byte {
	user := model.GetUser(username)
	blockedUser := UserToBlockReturned{
		Username: user.Username,
		Role:     user.Role,
		Blocked:  user.Blocked,
	}
	return createJson(blockedUser)
}

func retrieveOutOfWorkingHoursLogins(logins []model.Login, lastRun common_types.LastRun) []model.Login {
	var loginsToRetrieve []model.Login
	for _, login := range logins {
		if validators.ValidateLoginForIncidentResponse(login, lastRun) {
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

func decodeUserIncidentsBody(r *http.Request, structObject *common_types.LastRun) error {
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
