package api

import (
	"database/sql"
	"educations-castle/services/activity"
	"educations-castle/services/review"
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
	// router.Use(corsMiddleware) // Apply CORS middleware here
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.Use(corsMiddleware) // Apply CORS middleware here

	// User
	userCastle := user.NewCastle(s.db)
	userHandler := user.NewHandler(userCastle)
	userHandler.RegisterRoutes(subrouter)

	// Activity
	activityCastle := activity.NewCastle(s.db)
	activityHandler := activity.NewHandler(activityCastle, userCastle)
	activityHandler.RegisterRoutes(subrouter)

	// Review
	reviewCastle := review.NewCastle(s.db)
	reviewHandler := review.NewHandler(reviewCastle, userCastle)
	reviewHandler.RegisterRoutes(subrouter)

	log.Println(color.Format(color.GREEN, "Listening on "+s.addr))
	return http.ListenAndServe(s.addr, router)
}

// CORS middleware to add CORS headers to the HTTP response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
