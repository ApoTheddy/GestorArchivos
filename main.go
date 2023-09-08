package main

import (
	"fmt"
	"github/apotheddy/gestion-files-app/config"
	"github/apotheddy/gestion-files-app/database"
	"github/apotheddy/gestion-files-app/models"
	"github/apotheddy/gestion-files-app/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var DB database.Database

func main() {
	var dbconfig models.DatabaseConfigModel
	config.Getconfiguration(&dbconfig)
	DB.Init(dbconfig)
	DB.Connect()
	database.DB.AutoMigrate(models.FileModel{})

	isCreated, message := createDirUploads()
	if !isCreated {
		fmt.Println(message)
		return
	}
	r := mux.NewRouter()
	fileRoutes := routes.FileRoute{
		Route: r,
	}
	fileRoutes.Getroutes()
	uploadsHandler := http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads")))

	r.PathPrefix("/uploads/").Handler(uploadsHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Imágenes estáticas servidas con éxito!"))
	})

	http.ListenAndServe(":3000", r)
}

func createDirUploads() (isCreated bool, message string) {
	err := os.MkdirAll("uploads", os.ModePerm)
	isCreated = true
	if err != nil {
		message = "No se pudo crear correctamente la carpeta"
		isCreated = false
	}
	return
}
