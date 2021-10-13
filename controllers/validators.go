package controllers

import "fmt"

func validateUser(username, password string, user User) bool {
	return username == user.Username && password == user.Password && user.Blocked != true
}

func validateLoginHour(loginHour int) bool {
	fmt.Println(settings)
	return loginHour >= settings.From && loginHour <= settings.To
}

func validateLoginAttempt(username, password string) bool {
	for _, user := range users {
		fmt.Println(user)
		if validateUser(username, password, user){
			return true
		}
	}
	return false
}