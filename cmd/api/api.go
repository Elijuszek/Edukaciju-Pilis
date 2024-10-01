package api

import (
	"database/sql"
	"educations-castle/services/user"
	"educations-castle/utils/color"
	"log"
	"net/http"

	_ "educations-castle/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	userCastle := user.NewCastle(s.db)
	userHandler := user.NewHandler(userCastle)
	userHandler.RegisterRoutes(subrouter)
	log.Println(color.Format(color.GREEN, "Listening on "+s.addr))
	return http.ListenAndServe(s.addr, router)
}
