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
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type Handler struct {
	castle types.UserCastle
}

func NewHandler(castle types.UserCastle) *Handler {
	return &Handler{castle: castle}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users/login", h.handleLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/register", h.handleRegister).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/logout", auth.WithJWTAuth(h.handleLogout, h.castle)).Methods("POST", "OPTIONS")

	router.HandleFunc("/users", auth.WithJWTAuth(h.handleListUsers, h.castle, "administrator")).Methods("GET", "OPTIONS")

	router.HandleFunc("/users/{userID:[0-9]+}", auth.WithJWTAuth(h.handleGetUser, h.castle, "administrator", "organizer", "user")).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/update/{userID:[0-9]+}", auth.WithJWTAuth(h.handleUpdateUser, h.castle, "administrator", "organizer", "user")).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/delete/{userID:[0-9]+}", auth.WithJWTAuth(h.handleDeleteUser, h.castle, "administrator", "organizer", "user")).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/organizers/{organizerID:[0-9]+}", auth.WithJWTAuth(h.handleGetOrganizer, h.castle, "administrator", "organizer", "user")).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/create-organizer", auth.WithJWTAuth(h.handleCreateOrganizer, h.castle, "administrator")).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/create-administrator", auth.WithJWTAuth(h.handleCreateAdministrator, h.castle, "administrator")).Methods("POST", "OPTIONS")
}

// RegisterUser godoc
// @Summary      Create a new user account
// @Description  Create a new user by specifying the user information (username, email, password).
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        payload  body      types.UserPayload  true  "User registration data"
// @Success      201  {object}   types.UserResponse  "User successfully created"
// @Failure      400  {object}   types.ErrorResponse "invalud payload"
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

// LoginUser godoc
// @Summary      Login to user account
// @Description  Login to user account specifying (username, password).
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        payload  body   types.LoginUserPayload  true  "User login data"
// @Success      200  {object}   string  "Generated jwt"  "eyJhbGcifdghOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
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

	// TODO: GetUserByUsername and GetUserById also return hased password when it shouldn't
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
	role, err := resolveUserRole(u.ID, h.castle)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found user role"))
	}
	accessToken, err := auth.CreateJWT(secret, u.ID, role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	refreshToken, err := auth.CreateRefreshToken(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// TODO
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the request
	tokenString := utils.GetTokenFromRequest(r)

	// Validate the token
	token, err := auth.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid token, unable to logout"))
		return
	}

	// Get claims and extract the expiration time
	claims := token.Claims.(jwt.MapClaims)
	expiration, ok := claims["expiredAt"].(float64)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("couldn't parse token expiration"))
		return
	}

	// Add the token to the blacklist
	auth.AddTokenToBlacklist(tokenString, int64(expiration))

	// Respond with a successful logout message
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "successfully logged out"})
}

// ListUsers godoc
// @Summary      List all users
// @Description  List all registered users displaying the user information (username, email, password).
// @Tags         user
// @Produce      json
// @Success      200  {array}   types.User
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /users [get]
func (h *Handler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.castle.ListUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If no users found, return a message
	if len(reviews) == 0 {
		utils.WriteJSON(w, http.StatusOK, []types.Review{})
		return
	}

	// Return the users as a JSON response
	utils.WriteJSON(w, http.StatusOK, reviews)
}

// GetUser godoc
// @Summary      Get user by id
// @Description  Returns user with mathcing id
// @Tags         user
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200  {object}   types.User
// @Failure      400  {object}   types.ErrorResponse "missing user ID"
// @Failure      400  {object}   types.ErrorResponse "invalid user ID"
// @Failure      404  {object}   types.ErrorResponse "user not found"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /users/{userID} [get]
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

// GetOrganizer godoc
// @Summary      Get organizer by ID
// @Description  Returns organizer with matching ID
// @Tags         organizer
// @Produce      json
// @Param        organizerID  path      int  true  "Organizer ID"
// @Success      200  {object}   types.Organizer
// @Failure      400  {object}   types.ErrorResponse "missing organizer ID"
// @Failure      400  {object}   types.ErrorResponse "invalid organizer ID"
// @Failure      404  {object}   types.ErrorResponse "organizer not found"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /organizers/{organizerID} [get]
func (h *Handler) handleGetOrganizer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["organizerID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing organizer ID"))
		return
	}

	organizerID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid organizer ID"))
		return
	}

	// Retrieve the organizer by ID
	organizer, err := h.castle.GetOrganizerByID(organizerID)
	if err != nil || organizer == nil {
		if organizer == nil {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("organizer not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Return the organizer as a JSON response
	utils.WriteJSON(w, http.StatusOK, organizer)
}

// UpdateUser godoc
// @Summary      update user
// @Description  updates user with matching id with payload user data
// @Tags         user
// @Produce      json
// @Param        userID  path      int                 true  "User ID"
// @Param        payload body      types.UserPayload   true  "User update data"
// @Success      200     {object}  types.User
// @Failure      400     {object}  types.ErrorResponse "invalid payload"
// @Failure      404     {object}  types.ErrorResponse "user not found"
// @Failure      500     {object}  types.ErrorResponse "Internal server error"
// @Router       /users/update/{userID} [put]
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

	if !auth.CheckOwnership(r, existingUser.ID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
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

// DeleteUser godoc
// @Summary      delete user from database
// @Description  deletes user with specified id
// @Tags         user
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200  {object}   types.ErrorResponse "user with id %d successfully deleted"
// @NoContent    204  {object}   types.ErrorResponse "user not found"
// @Failure      400  {object}   types.ErrorResponse "invalid payload"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /users/delete/{userID} [delete]
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
			utils.WriteError(w, http.StatusNoContent, fmt.Errorf("user not found"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching review: %w", err))
		return
	}

	if !auth.CheckOwnership(r, existingUser.ID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
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

// CreateOrganizer godoc
// @Summary      create organizer role inside database
// @Description  creates organizer role inside database with specified description
// @Tags         user
// @Produce      json
// @Param        payload  body   types.CreateOrganizerPayload  true  "Organizer data"
// @Success      201  {object}   types.ErrorResponse "Organizer with ID %d successfully created"
// @Failure      400  {object}   types.ErrorResponse "invalid payload"
// @Failure      404  {object}   types.ErrorResponse "user not found"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /users/create-organizer [POST]
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
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user %d not found", payload.ID))
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

// TODO:
func (h *Handler) handleCreateAdministrator(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.CreateAdministratorPayload
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
	// TODO: check if admin already exists
	_, err := h.castle.GetUserByID(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user %d not found", payload.ID))
		return
	}

	// if it doesnt  create the new user
	err = h.castle.CreateAdministrator(types.CreateAdministratorPayload{
		ID:            payload.ID,
		SecurityLevel: payload.SecurityLevel,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Organizer with ID %d successfully created", payload.ID))
}

func resolveUserRole(userID int, castle types.UserCastle) (string, error) {
	// Check if the user is an administrator
	if admin, err := castle.GetAdministratorByID(userID); err == nil && admin != nil {
		return "administrator", nil
	}

	// Check if the user is an organizer
	if organizer, err := castle.GetOrganizerByID(userID); err == nil && organizer != nil {
		return "organizer", nil
	}

	// Check if the user is a regular user
	if user, err := castle.GetUserByID(userID); err == nil && user != nil {
		return "user", nil
	}

	return "", fmt.Errorf("user not found in any role")
}
