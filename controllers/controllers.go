package controllers

func InitControllers() {
	loadModelToMemoryJson(&settings, &users, &logins)
	initApi()
}
