package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	Title string
	Body  []byte
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

func parseSettingsRequestForm(r *http.Request) (from int, to int) {
	r.ParseForm()
	fmt.Println(r.Form)
	for key, value := range r.Form {
		if key == "from" {
			from = populateKey(from, value[0])
		}
		if key == "to" {
			to = populateKey(from, value[0])
		}
		fmt.Printf("%s = %s\n", key, value)
	}
	return from, to
}

func populateKey(key int, value string) int {
	i, err := strconv.Atoi(value)
	handleParsingError(err)
	key = i
	return key
}

func handleParsingError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getLoginAttempt(username string, r *http.Request) Login {
	user := getUserByUsername(username)
	loginHour := time.Now().Hour()
	return Login{
		LoginId: rand.Intn(1000000),
		Username: user.Username,
		UserEmail: user.Email,
		Ip: getIP(r),
		UserRole: user.Role,
		LoginDate: time.Now().Format("2006-01-02 15:04:05"),
		LoginHour: loginHour,
		ValidLoginHour: validateLoginHour(loginHour),
	}
}

func getUserByUsername(username string) User{
	var retutnUser User
	for _, user := range users {
		fmt.Println(user)
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
	cookie := http.Cookie{Name: "token", Value: "authenticated", Expires:expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return w
}

func deleteAuthCookie(w http.ResponseWriter) http.ResponseWriter {
	expiration := time.Unix(0, 0)
	cookie := http.Cookie{Name: "token", Value: "", Expires:expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return w
}