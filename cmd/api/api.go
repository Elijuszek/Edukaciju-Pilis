package api

import (
	"database/sql"
	"educations-castle/service/user"
	"educations-castle/utils/color"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)
	log.Println(color.Format(color.GREEN, "Listening to "+s.addr))

	return http.ListenAndServe(s.addr, router)
}