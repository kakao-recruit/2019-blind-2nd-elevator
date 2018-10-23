package main

import (
	_ "2019-blind-2nd-elevator/elevator/config"
	"github.com/gorilla/mux"
	"net/http"

	"2019-blind-2nd-elevator/elevator/internal/api"
	"2019-blind-2nd-elevator/elevator/internal/myauth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/start/{userid}/{problem}/{count}", api.Start).Methods("POST")
	router.HandleFunc("/oncalls", api.OnCalls).Methods("GET")
	router.HandleFunc("/action", api.Action).Methods("POST")

	// For viewer
	router.HandleFunc("/viewer", Index).Methods("GET")
	router.HandleFunc("/viewer/{filename}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "viewer/" + mux.Vars(r)["filename"])
	}).Methods("GET")
	router.HandleFunc("/viewer/trials/{token}", Trials).Methods("GET")

	amw := myauth.AuthMiddleware{}

	router.Use(amw.Middleware)

	log.Debugf("Ready: %s", viper.GetString("ListenAddr"))
	log.Fatal(http.ListenAndServe(viper.GetString("ListenAddr"), router))
}
