package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"educations-castle/cmd/api"
	"educations-castle/configs"
	"educations-castle/db"
	"educations-castle/utils/color"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// Load the CA certificate from the path specified in the configuration
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(configs.Envs.CACertPath) // Assuming CACertPath is the config field
	if err != nil {
		log.Fatalf("Failed to read CA certificate from path %s: %v", configs.Envs.CACertPath, err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append CA certificate")
	}

	// Register a custom tls.Config
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})

	// Use the custom TLS config in the DB connection
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		TLSConfig:            "custom", // This references the TLS config we registered
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
