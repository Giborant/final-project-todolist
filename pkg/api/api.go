package api

import (
	"net/http"
)

const (
	webDir = "web"
)

func Init() {
	// r := mux.NewRouter()
	// r.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
}
