package services

import (
	"encoding/json"
	"fmt"
	"github/apotheddy/gestion-files-app/database"
	"github/apotheddy/gestion-files-app/models"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func errorRequest(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, message)))
}

func Getfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var files []models.FileModel
	database.DB.Where("dir = ''").Find(&files)
	json.NewEncoder(w).Encode(files)
}

func GetOnlyPhotosAndVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var files []models.FileModel
	database.DB.Where("RIGHT(filename,4) IN ('.jpg','.png','.jpeg','.mkv','.mp4')").Find(&files)
	json.NewEncoder(w).Encode(files)
}

func Getfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var file models.FileModel
	queryWhere := "ID = " + params["id"]
	database.DB.Where(queryWhere).First(&file)

	if file.ID == 0 {
		errorRequest(w, "El archivo no se ha encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(file)
	return
}

func GetDirs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dirs, err := os.ReadDir("uploads/")

	if err != nil {
		errorRequest(w, "No se pudieron leer los directorios", http.StatusInternalServerError)
		return
	}

	var listDirs []string

	for _, dir := range dirs {
		if dir.IsDir() {
			listDirs = append(listDirs, dir.Name())
		}
	}
	if listDirs != nil {
		json.NewEncoder(w).Encode(listDirs)
		return
	}
	w.Write([]byte("[]"))
}

type Dir struct {
	Dirname string `json:"dirname"`
}

func CreateDir(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dir Dir
	json.NewDecoder(r.Body).Decode(&dir)

	if dir.Dirname == "" {
		errorRequest(w, "Ingrese un nombre para la carpeta a crear", http.StatusBadRequest)
		return
	}
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		errorRequest(w, fmt.Sprintf("Error al crear la carpeta de %s", dir), http.StatusInternalServerError)
		return
	}
	err = os.MkdirAll(fmt.Sprintf("uploads/%s/", dir.Dirname), os.ModePerm)
	if err != nil {
		errorRequest(w, fmt.Sprintf("No se pudo crear la carpeta '%s'", dir), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"message":"Carpeta creada correctamente"}`))
}

func GetFilesByDirName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var files []models.FileModel
	params := mux.Vars(r)
	queryWhere := fmt.Sprintf("dir = '%s'", params["dirname"])
	database.DB.Where(queryWhere).Find(&files)
	json.NewEncoder(w).Encode(files)
}

func RemoveFileById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	var file models.FileModel
	database.DB.First(&file, id)
	if file.Filename == "" {
		errorRequest(w, "El archivo a eliminar no existe", http.StatusNotFound)
		return
	}
	err := os.Remove(file.FilePath)
	if err != nil {
		errorRequest(w, "No se pudo eliminar el archivo", http.StatusInternalServerError)
		return
	}
	database.DB.Unscoped().Delete(&file)
	w.Write([]byte(`{"message":"Archivo eliminado correctamente"}`))
}

func RemoveDir(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dir Dir
	var files []models.FileModel
	json.NewDecoder(r.Body).Decode(&dir)
	path := "uploads/" + dir.Dirname
	err := os.RemoveAll(path)
	if err != nil {
		errorRequest(w, "No se pudo eliminar la carepta", http.StatusInternalServerError)
		return
	}
	queryWhere := "file_path LIKE 'uploads/" + dir.Dirname + "/%'"
	database.DB.Where(queryWhere).Find(&files)
	database.DB.Unscoped().Delete(files)
	w.Write([]byte(`{"message":"Carpeta eliminada correctamente"}`))
}

func Uploadfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, handler, err := r.FormFile("file")
	dirname := r.FormValue("dirname")

	if err != nil {
		errorRequest(w, "Indique el archivo que desea subir", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filepath := "uploads/"
	if dirname != "" {
		filepath += dirname + "/"
	}
	uniquename := fmt.Sprintf("%s%s", uuid.New().String(), path.Ext(handler.Filename))

	msg, er := createFilesFolder(w, uniquename, filepath, file)

	if er != nil {
		errorRequest(w, msg, http.StatusInternalServerError)
		return
	}

	_, newFile, err := saveDataFileInDatabase(handler.Filename, uniquename, filepath, dirname)

	if err != nil {
		errorRequest(w, "Indique el archivo que desea subir", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newFile)
}

func saveDataFileInDatabase(filename, uniquename, filepath, dirname string) (message string, file models.FileModel, err error) {

	file = models.FileModel{
		Filename:   filename,
		FilePath:   filepath + uniquename,
		Dir:        dirname,
		UniqueName: uniquename,
	}
	newFile := database.DB.Create(&file)

	if err := newFile.Error; err != nil {
		message = "No se pudo guardar el archivo en la base de datos"
		err = fmt.Errorf("No se guardo correctamente")
	} else {
		message = "Archivo guardado correctamente"
		err = nil
	}
	return message, file, err
}

func createFilesFolder(w http.ResponseWriter, filename string, filepath string, file io.Reader) (msg string, er error) {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		msg = "No se pudo crear la carpeta de 'uploads'"
		er = fmt.Errorf(msg)
	}
	out, err := os.Create(filepath + filename)
	defer out.Close()
	if err != nil {
		msg = "No se pudo abrir el archivo correctamente"
		er = fmt.Errorf(msg)
	}
	_, err = io.Copy(out, file)
	if err != nil {
		msg = "Ocurrio un problema al guardar la informacion del archivo"
	}
	return
}
