package validators

import (
	"fmt"
	"github.com/Victor-Yurievich/working-hours-app/common-types"
	"github.com/Victor-Yurievich/working-hours-app/model"
	"os"
	"time"
)

//type LastRun struct {
//	Time string `json:"time"`
//}

//Ask Lior about common_types interoperability

func validateUser(username, password string, user model.User) bool {
	return username == user.Username && password == user.Password && user.Blocked != true
}

func ValidateLoginHour(loginHour int) bool {
	return loginHour >= model.TimeSettings.From && loginHour <= model.TimeSettings.To
}

func ValidateLoginAttempt(username, password string) bool {
	for _, user := range model.Users {
		if validateUser(username, password, user) {
			return true
		}
	}
	return false
}

//func ValidateLoginForIncidentResponse(login model.Login, lastRun LastRun) bool {
//	return login.ValidLoginHour == false && inTimeInterval(login.LoginDate, lastRun)
//}

//Ask Lior about common types interoperability

func ValidateLoginForIncidentResponse(login model.Login, lastRun common_types.LastRun) bool {
	return login.ValidLoginHour == false && inTimeInterval(login.LoginDate, lastRun)
}

//Ask Lior / Shani how we do generic functions
func inTimeInterval(loginDate string, lastRun common_types.LastRun) bool {
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
