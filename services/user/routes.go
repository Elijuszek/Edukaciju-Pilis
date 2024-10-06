package user

import (
	"educations-castle/configs"
	"educations-castle/services/auth"
	"educations-castle/types"
	"educations-castle/utils"
	"fmt"
	"net/http"

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
	router.HandleFunc("/login", h.handleLogin).Methods(("POST"))
	router.HandleFunc("/register", h.handleRegister).Methods(("POST"))

	router.HandleFunc("/delete-user", h.handleDeleteUser).Methods(("POST"))
	router.HandleFunc("/create-organizer", h.handleCreateOrganizer).Methods(("POST"))
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

	//utils.WriteJSON(w, http.StatusOK, map[string]string{"token": ""})
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
	var payload types.RegisterUserPayload
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

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload struct {
		ID int `json:"id" validate:"required"`
	}

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

	// delete organizer if exists
	err := h.castle.DeleteOrganizer(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete organizer: %v", err))
		return
	}

	// delete administrator if exists
	err = h.castle.DeleteAdministrator(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete administrator: %v", err))
		return
	}

	// delete user itself
	err = h.castle.DeleteUser(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete user: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("user with name %d successfully deleted", payload.ID))
}

func (h *Handler) handleCreateOrganizer(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.CreateOrganizerPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if the user exists
	_, err := h.castle.GetUserByID(payload.ID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with id %d doesn't exists", payload.ID))
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

	utils.WriteJSON(w, http.StatusCreated, nil)
}
