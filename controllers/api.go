package controllers

import (
	"log"
	"net/http"
)

func initApi() {
	http.HandleFunc("/", handler) // tells the http package to handle all requests to the web root ("/") with handler.
	http.HandleFunc("/login/", makeHandler(loginHandler))
	http.HandleFunc("/auth/", makeHandler(loginAuthHandler))
	http.HandleFunc("/dashboard/", makeHandler(authHandler))
	http.HandleFunc("/settings/", makeHandler(settingsHandler))
	http.HandleFunc("/ws/", wsEndpoint)
	log.Fatal(http.ListenAndServe(":8088", nil))
}
