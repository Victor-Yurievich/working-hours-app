package main

import (
	"github.com/Victor-Yurievich/working-hours-app/api"
	"github.com/Victor-Yurievich/working-hours-app/model"
)

func main() {
	model.LoadModelToMemoryJson(&model.TimeSettings, &model.Users, &model.Logins)
	api.InitApi()
}
