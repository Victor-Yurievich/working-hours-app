package controllers

func InitControllers() {
	LoadModelToMemoryJson(&settings, &users, &logins)
	initApi()
}
