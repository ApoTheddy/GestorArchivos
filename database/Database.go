package database

import (
	"fmt"
	"github/apotheddy/gestion-files-app/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DatabaseMethods interface {
	Init(Config models.DatabaseConfigModel) *Database
	Connect()
}

var DB *gorm.DB

type Database struct {
	Config models.DatabaseConfigModel
}

func (d *Database) Init(config models.DatabaseConfigModel) {
	d.Config = config
}

func (d Database) Connect() {
	var err error
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", d.Config.User, d.Config.Password, d.Config.Host, d.Config.Port, d.Config.DatabaseName)
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Ocurrio un error al conectar a la base de datos: ", err)
		return
	}
	fmt.Println("Conectado correctamente a la base de datos")
}
