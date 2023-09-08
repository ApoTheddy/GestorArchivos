package routes

import (
	"github/apotheddy/gestion-files-app/controllers"

	"github.com/gorilla/mux"
)

type FileRouteMethods interface {
	Getroutes()
}

type FileRoute struct {
	Route *mux.Router
}

func (fr FileRoute) Getroutes() {
	fr.Route.HandleFunc("/files", controllers.Getfiles).Methods("GET")
	fr.Route.HandleFunc("/files/{id}", controllers.Getfile).Methods("GET")
	fr.Route.HandleFunc("/photos", controllers.GetOnlyPhotosAndVideos).Methods("GET")
	fr.Route.HandleFunc("/files/search/{dirname}", controllers.GetFilesByDirName).Methods("GET")
	fr.Route.HandleFunc("/dirs", controllers.GetDirs).Methods("GET")

	fr.Route.HandleFunc("/files", controllers.UploadFile).Methods("POST")
	fr.Route.HandleFunc("/files/dir", controllers.CreateDir).Methods("POST")
	fr.Route.HandleFunc("/files/{id}", controllers.RemoveFileById).Methods("DELETE")
	fr.Route.HandleFunc("/files/dir/delete", controllers.RemoveDir).Methods("DELETE")
}
