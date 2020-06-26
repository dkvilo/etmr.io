package main

import (
	"log"
	"net/http"
	controller "github.com/dkvilo/etmr.io/controllers"
	"github.com/dkvilo/etmr.io/core"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
)

func init() {
	err := env.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
	}
}

func main()  {

	// Crate New App
	app := core.App{}
	defer app.MongoClient.Disconnect(app.DatabaseConnectionContext)

	// Create Mux Router
	router := mux.NewRouter()

	// Establish Database Connection
	app.ConnectToMongoDB(true)

	// Pass App instance to the API Controller
	api := controller.APIController {
			controller.BaseController { App: &app },
			nil,
	}

	// Host static folder
	go router.Handle("/", http.FileServer(http.Dir("./public")))

	// URL API EP(S)
	router.HandleFunc("/api/url", api.URL).Methods("POST")

	// Redirect Root Exit Node
	router.HandleFunc("/{slug}", api.RedirectURL).Methods("GET")

	http.ListenAndServe(":3000", router)
}


