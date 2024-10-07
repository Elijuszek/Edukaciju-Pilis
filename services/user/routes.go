package user

import (
	"database/sql"
	"educations-castle/configs"
	"educations-castle/services/auth"
	"educations-castle/types"
	"educations-castle/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	castle types.UserCastle
}

func NewHandler(castle types.UserCastle) *Handler {
	return &Handler{castle: castle}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users/login", h.handleLogin).Methods(("POST"))
	router.HandleFunc("/users/register", h.handleRegister).Methods(("POST"))

	router.HandleFunc("/users", h.handleListUsers).Methods(("GET"))
	router.HandleFunc("/users/{userID:[0-9]+}", h.handleGetUser).Methods(("GET"))
	router.HandleFunc("/users/update/{userID:[0-9]+}", h.handleUpdateUser).Methods("PUT")
	router.HandleFunc("/users/delete/{userID:[0-9]+}", h.handleDeleteUser).Methods(("DELETE"))

	router.HandleFunc("/users/create-organizer", h.handleCreateOrganizer).Methods(("POST"))
}

// LoginUser godoc
// @Summary      Login to user account
// @Description  Login to user account specifying (username, password).
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        payload  body   types.LoginUserPayload  true  "User login data"
// @Success      201  {object}   jwt.token  "User successfully loged in"
// @Failure      400  {object}   types.ErrorResponse "Invalid payload"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /users/login [post]
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.castle.GetUserByUsername(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found invalid username or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found invalid username or password"))
		return
	}

	// JWT
	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// RegisterUser godoc
// @Summary      Create a new user account
// @Description  Create a new user by specifying the user information (username, email, password).
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload  body      types.RegisterUserPayload  true  "User registration data"
// @Success      201  {object}   types.UserResponse  "User successfully created"
// @Failure      400  {object}   types.ErrorResponse "Invalid payload or user already exists"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /users/register [post]
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.UserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if the user exists
	_, err := h.castle.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	_, err = h.castle.GetUserByUsername(payload.Username)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with username %s already exists", payload.Username))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it doesnt  create the new user
	err = h.castle.CreateUser(types.User{
		Username: payload.Username,
		Password: hashedPassword,
		Email:    payload.Email,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("User %s successfully registered", payload.Username))
}

func (h *Handler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.castle.ListUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If no users found, return a message
	if len(reviews) == 0 {
		utils.WriteJSON(w, http.StatusOK, "no reviews found")
		return
	}

	// Return the users as a JSON response
	utils.WriteJSON(w, http.StatusOK, reviews)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.castle.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the review ID from the URL parameters
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	// Convert review ID from string to int
	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	// Get JSON payload
	var payload types.UserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// Check if the user exists
	existingUser, err := h.castle.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	_, err = h.castle.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with new email %s already exists", payload.Email))
		return
	}

	_, err = h.castle.GetUserByUsername(payload.Username)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with new username %s already exists", payload.Username))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Update the user
	updateUser := types.User{
		ID:               existingUser.ID,
		Username:         payload.Username,
		Password:         hashedPassword,
		Email:            payload.Email,
		RegistrationDate: existingUser.RegistrationDate,
		LastLoginDate:    existingUser.LastLoginDate,
	}

	err = h.castle.UpdateUser(updateUser)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updateUser)
}

func (h *Handler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	// Check if the user exists
	existingUser, err := h.castle.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching review: %w", err))
		return
	}

	// delete user
	err = h.castle.DeleteUser(existingUser.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete user: %d", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("user with id %d successfully deleted", userID))
}

func (h *Handler) handleCreateOrganizer(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.CreateOrganizerPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if the user exists
	// TODO: check if organizer already exists
	_, err := h.castle.GetUserByID(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with id %d already exists", payload.ID))
		return
	}

	// if it doesnt  create the new user
	err = h.castle.CreateOrganizer(types.Organizer{
		ID:          payload.ID,
		Description: &payload.Description,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Organizer with ID %d successfully created", payload.ID))
}
