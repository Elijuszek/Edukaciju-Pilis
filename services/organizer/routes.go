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

	router.HandleFunc("/create-package", h.handleCreatePackage).Methods(("POST"))

	router.HandleFunc("/register-organizer", h.regi).Methods(("POST"))
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

	err := h.castle.CreateActivity(types.Activity{
		Name:         payload.Name,
		Description:  payload.Description,
		BasePrice:    payload.BasePrice,
		Hidden:       payload.Hidden,
		Category:     payload.Category,
		Fk_PackageId: payload.Fk_PackageId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
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
	var payload types.PackagePayload
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
		Name:           payload.Name,
		Description:    payload.Description,
		Price:          payload.Price,
		Fk_OrganizerId: payload.Fk_OrganizerId,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
