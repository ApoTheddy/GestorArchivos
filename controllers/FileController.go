package controllers

import (
	"github/apotheddy/gestion-files-app/services"
	"net/http"
)

func Getfiles(w http.ResponseWriter, r *http.Request) {
	services.Getfiles(w, r)
}
func GetOnlyPhotosAndVideos(w http.ResponseWriter, r *http.Request) {
	services.GetOnlyPhotosAndVideos(w, r)
}

func Getfile(w http.ResponseWriter, r *http.Request) {
	services.Getfile(w, r)
}

func GetDirs(w http.ResponseWriter, r *http.Request) {
	services.GetDirs(w, r)
}

func GetFilesByDirName(w http.ResponseWriter, r *http.Request) {
	services.GetFilesByDirName(w, r)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	services.Uploadfile(w, r)
}

func CreateDir(w http.ResponseWriter, r *http.Request) {
	services.CreateDir(w, r)
}

func RemoveFileById(w http.ResponseWriter, r *http.Request) {
	services.RemoveFileById(w, r)
}

func RemoveDir(w http.ResponseWriter, r *http.Request) {
	services.RemoveDir(w, r)
}
