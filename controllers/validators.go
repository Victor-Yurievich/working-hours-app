package controllers

import (
	"fmt"
	"os"
	"time"
)

func validateUser(username, password string, user User) bool {
	return username == user.Username && password == user.Password && user.Blocked != true
}

func validateLoginHour(loginHour int) bool {
	return loginHour >= settings.From && loginHour <= settings.To
}

func validateLoginAttempt(username, password string) bool {
	for _, user := range users {
		if validateUser(username, password, user) {
			return true
		}
	}
	return false
}

func validateLoginForIncidentResponse(login Login, lastRun LastRun) bool {
	return login.ValidLoginHour == false && inTimeInterval(login.LoginDate, lastRun)
}

func inTimeInterval(loginDate string, lastRun LastRun) bool {
	endDate, err := time.Parse("2006-01-02T15:04:05.000Z", time.Now().Format("2006-01-02T15:04:05.000Z"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	to := endDate.Unix()
	startDate, err := time.Parse("2006-01-02T15:04:05.000Z", lastRun.Time)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	from := startDate.Unix()
	dateToCheck, err := time.Parse("2006-01-02T15:04:05.000Z", loginDate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	timeStampToCheck := dateToCheck.Unix()
	return from < timeStampToCheck && timeStampToCheck < to
}
