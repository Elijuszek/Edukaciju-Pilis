package activity

// TODO: Use sqlx
import (
	"educations-castle/types"
	"educations-castle/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	castle types.ActivityCastle
}

func NewHandler(castle types.ActivityCastle) *Handler {
	return &Handler{castle: castle}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create-activity", h.handleCreateActivity).Methods(("POST"))
	router.HandleFunc("/delete-activity", h.handleDeleteActivity).Methods(("DELETE"))

	router.HandleFunc("/create-package", h.handleCreatePackage).Methods(("POST"))
	router.HandleFunc("/delete-package", h.handleDeletePackage).Methods(("DELETE"))

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
func (h *Handler) handleCreateActivity(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.CreateActivityPayload
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

	// check if the activity exists inside package
	_, err := h.castle.GetActivityInsidePackageByName(payload.Name, payload.FkPackageID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("activity with name %s inside package already exists", payload.Name))
		return
	}

	err = h.castle.CreateActivity(types.Activity{
		Name:        payload.Name,
		Description: payload.Description,
		BasePrice:   payload.BasePrice,
		Hidden:      payload.Hidden,
		Category:    payload.Category,
		FkPackageID: payload.FkPackageID,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleDeleteActivity(w http.ResponseWriter, r *http.Request) {
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

	// delete the activity
	err := h.castle.DeleteActivity(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete activity: %v", err))
		return
	}

	// TODO: status NO content
	utils.WriteJSON(w, http.StatusAccepted, fmt.Sprintf("activity with name %d successfully deleted", payload.ID))
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
func (h *Handler) handleCreatePackage(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.CreatePackagePayload
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

	// check if the package exists
	_, err := h.castle.GetPackageByName(payload.Name)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("package with name %s already exists", payload.Name))
		return
	}

	// if not create
	err = h.castle.CreatePackage(types.Package{
		Name:          payload.Name,
		Description:   payload.Description,
		Price:         payload.Price,
		FkOrganizerID: payload.FkOrganizerID,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleDeletePackage(w http.ResponseWriter, r *http.Request) {
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

	// delete package
	err := h.castle.DeletePackage(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete package: %v", err))
		return
	}

	// TODO: status NO content
	utils.WriteJSON(w, http.StatusAccepted, fmt.Sprintf("package with id %d successfully deleted", payload.ID))
}
