package main

import (
	"gt-alann/config"
	apimanagement "gt-alann/internal/apiManagement"
	"gt-alann/internal/handlers"
	"gt-alann/internal/serverManagement"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", handlers.AllArtistsHandle)
	http.HandleFunc("/artist", handlers.ArtistHandle)
	http.HandleFunc("/admin", handlers.AdminHandle)
	appConfig := config.ConfigLoad()

	// config
	if templateCache, err := handlers.CreateTemplateCache(); err != nil {
		log.Fatal(err)
	} else {
		appConfig.TemplateCache = templateCache
		appConfig.Port = ":8080"
		appConfig.Api = "https://groupietrackers.herokuapp.com/api"
	}
	apimanagement.ConfigApi()
	handlers.ConfigHandle()

	log.Println("Server start on http://localhost" + appConfig.Port + "/")
	serveur := http.Server{Addr: appConfig.Port}
	go serveur.ListenAndServe()
	serverManagement.WaitServerOrder(&serveur)
}
