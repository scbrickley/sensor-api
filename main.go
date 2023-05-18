package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scbrickley/sensor-api/dependencies"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	deps, err := dependencies.NewDependencies()
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/sensors", deps.ListSensorsHandler()).Methods("GET")
	router.HandleFunc("/sensors", deps.InsertSensorHandler()).Methods("POST")
	router.HandleFunc("/sensors/nearest", deps.NearestSensorHandler()).Methods("GET")
	router.HandleFunc("/sensors/{name}", deps.GetSensorByNameHandler()).Methods("GET")
	router.HandleFunc("/sensors/{name}", deps.UpdateSensorHandler()).Methods("PUT")
	router.HandleFunc("/sensors/{name}", deps.DeleteSensorHandler()).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
