package main

import (
	"escapade/internal/config"
	"escapade/internal/database"
	mi "escapade/internal/middleware"
	"escapade/internal/services/api"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	confPath = "conf.json"
)

func main() {
	fmt.Println("Ok")
	conf, confErr := config.Init(confPath)
	fmt.Println("Ok")
	if confErr != nil {
		panic(confErr)
	}
	fmt.Println("Ok")
	db, dbErr := database.Init()
	if dbErr != nil {
		panic(dbErr)
	}
	fmt.Println("Ok")
	API := api.Init(db)

	r := mux.NewRouter()

	r.HandleFunc("/", mi.CORS(conf.Cors)(API.Ok))
	r.HandleFunc("/register", mi.CORS(conf.Cors)(API.Register)).Methods("POST")
	r.HandleFunc("/delete", mi.CORS(conf.Cors)(API.DeleteAccount)).Methods("DELETE")
	r.HandleFunc("/delete", mi.CORS(conf.Cors)(API.DeleteAccountOptions)).Methods("OPTIONS")

	r.HandleFunc("/login", mi.CORS(conf.Cors)(API.Login)).Methods("POST")
	r.HandleFunc("/me", mi.CORS(conf.Cors)(API.Me)).Methods("GET")
	r.HandleFunc("/{name}/games", mi.CORS(conf.Cors)(API.GetPlayerGames)).Methods("GET")
	r.HandleFunc("/{name}/profile", mi.CORS(conf.Cors)(API.GetProfile)).Methods("GET")

	fmt.Println("launched, look at us on " + conf.Server.Host + conf.Server.Port) //+ os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), r)
	//err := http.ListenAndServe(conf.Server.Port, r)
	fmt.Println("oh, this is error:" + err.Error())
}
