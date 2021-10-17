package api

import (
	"fmt"
	"github.com/Victor-Yurievich/working-hours-app/websockets"
	"log"
	"net/http" //Ask Lior about Router
)

func InitApi() {
	http.HandleFunc("/", handler) // tells the http package to handle all requests to the web root ("/") with handler.
	http.HandleFunc("/login", makeHandler(loginHandler))
	http.HandleFunc("/logout", makeHandler(logoutHandler))
	http.HandleFunc("/auth", makeHandler(loginAuthHandler))
	http.HandleFunc("/dashboard", makeHandler(authHandler))
	http.HandleFunc("/settings", makeHandler(settingsHandler))
	http.HandleFunc("/ws", makeHandler(websockets.WsEndpoint))
	http.HandleFunc("/block-user", makeHandler(blockUserHandler))
	http.HandleFunc("/ping", makeHandler(pongResponse))
	http.HandleFunc("/fetch-incidents", makeHandler(incidentsResponse))
	fmt.Println("Server is listening on the port 8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
