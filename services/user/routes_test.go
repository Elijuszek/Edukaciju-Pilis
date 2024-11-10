package user

import (
	"bytes"
	"educations-castle/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserServiceHandler(t *testing.T) {

	userCastle := &mockUserCastle{}
	handler := NewHandler(userCastle)

	t.Run("Should fail if the user payload is invalid", func(*testing.T) {
		payload := types.UserPayload{
			Username: "user",
			Password: "password",
			Email:    "invalid",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should correctly register the user", func(t *testing.T) {
		payload := types.UserPayload{
			Username: "user",
			Password: "password",
			Email:    "validUser@email.com",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserCastle struct{}

func (m *mockUserCastle) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (c *mockUserCastle) GetUserByUsername(username string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserCastle) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserCastle) CreateUser(types.User) error {
	return nil
}

func (m *mockUserCastle) UpdateUser(types.User) error {
	return nil
}

func (m *mockUserCastle) DeleteUser(id int) error {
	return nil
}

func (m *mockUserCastle) ListUsers() ([]*types.User, error) {
	return nil, nil
}

func (m *mockUserCastle) CreateOrganizer(types.Organizer) error {
	return nil
}

func (m *mockUserCastle) GetAdministratorByID(id int) (*types.Administrator, error) {
	return nil, nil
}

func (m *mockUserCastle) GetOrganizerByID(id int) (*types.Organizer, error) {
	return nil, nil
}

func (m *mockUserCastle) CreateAdministrator(types.Administrator) error {
	return nil
}
