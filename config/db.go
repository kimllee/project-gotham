package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Datasource Name
	dsn := "api_partage:api_partage@tcp(127.0.0.1:3306)/partage_bdd?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Erreur de connexion à la base de données :", err)
	}

	DB = db
}
