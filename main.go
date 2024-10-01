package main

import (
	"database/sql"
	"educations-castle/cmd/api"
	"educations-castle/configs"
	"educations-castle/db"
	"educations-castle/utils/color"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(color.Format(color.GREEN, "DB: Successfully connected!"))
}
