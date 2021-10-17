package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Settings struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Blocked  bool   `json:"blocked"`
}

type Login struct {
	LoginId        string `json:"login_id"`
	Username       string `json:"username"`
	UserEmail      string `json:"user_email"`
	Ip             string `json:"ip"`
	UserRole       string `json:"user_role"`
	LoginDate      string `json:"login_date"`
	LoginHour      int    `json:"login_hour"`
	ValidLoginHour bool   `json:"valid_login_hour"`
}

var TimeSettings = Settings{}
var Users = []User{}
var Logins = []Login{}

func LoadModelToMemoryJson(TimeSettings *Settings, users *[]User, logins *[]Login) {
	loadSettings(TimeSettings)
	loadUsers(users)
	loadLogins(logins)
}

func loadFile(name string, extension string) ([]byte, error) {
	filename := "./model/file-db/" + name + "." + extension
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func loadSettings(settings *Settings) { // Ask Lior about generic typing
	loadedFile, err := loadFile("settings", "json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(loadedFile, &settings)
}

func loadUsers(users *[]User) {
	loadedFile, err := loadFile("users", "json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(loadedFile, &users)
}

func loadLogins(logins *[]Login) {
	loadedFile, err := loadFile("logins", "json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(loadedFile, &logins)
}

func SaveLogin(logins *[]Login) error {
	filename := "./model/file-db/logins.json"
	loginsJson, err := json.Marshal(logins)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(filename, loginsJson, 0600)
}

func SaveSettings(settings *Settings) error {
	filename := "./model/file-db/settings.json"
	settingsJson, err := json.Marshal(settings)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(filename, settingsJson, 0600)
}

func SaveUsers() error {
	filename := "./model/file-db/users.json"
	usersJson, err := json.Marshal(Users)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(filename, usersJson, 0600)
}

func GetUser(username string) User {
	var userToReturn = User{}
	for _, user := range Users {
		if user.Username == username {
			userToReturn = user
		}
	}
	return userToReturn
}
