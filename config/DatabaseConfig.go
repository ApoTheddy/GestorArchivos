package config

import (
	"fmt"
	"github/apotheddy/gestion-files-app/models"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Getconfiguration(config *models.DatabaseConfigModel) {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error cargando el archivo .env", err)
		return
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		fmt.Println("El puerto obtenido del '.env' no corresponde a un valor numerico")
		return
	}

	*config = models.DatabaseConfigModel{
		Host:         os.Getenv("HOST"),
		Port:         port,
		User:         os.Getenv("USER"),
		Password:     os.Getenv("SA_PASSWORD"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
	}
}
