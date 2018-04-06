package main

import (
	"api-barebone/handlers"
	authHandler "api-barebone/handlers/bbAuth"
	"api-barebone/src/bbAuth"
	"log"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup the authentication handler
	authHandler.AuthHandler = &bbAuth.PlainTextAuthHandler{
		AuthFile:    "./.secrets.csv",
		SecretsFile: "./.secrets.json",
	}

	// Setup the router
	router := vestigo.NewRouter()
	vestigo.AllowTrace = false
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      []string{"*"},
		AllowCredentials: true,
	})

	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		handler := handlers.Routes.GetAction(r.Method, r.RequestURI)

		if handler != nil {
			handler(w, r)
		}
	})

	logrus.WithFields(logrus.Fields{
		"port": 8000,
	}).Info("Initializing server")
	log.Fatal(http.ListenAndServe(":8000", router))
}
