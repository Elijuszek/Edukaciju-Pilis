package activity

// TODO: Use sqlx
import (
	"database/sql"
	"educations-castle/types"
	"educations-castle/utils"
	"fmt"
	"net/http"
	"strconv"

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
	router.HandleFunc("/activities", h.handleListActivities).Methods(("GET"))
	router.HandleFunc("/activities/create", h.handleCreateActivity).Methods(("POST"))
	router.HandleFunc("/activities/{activityID:[0-9]+}", h.handleGetActivity).Methods(("GET"))
	router.HandleFunc("/activities/update/{activityID:[0-9]+}", h.handleUpdateActivity).Methods(("PUT"))
	router.HandleFunc("/activities/delete/{activityID:[0-9]+}", h.handleDeleteActivity).Methods(("DELETE"))

	router.HandleFunc("/packages/create", h.handleCreatePackage).Methods(("POST"))
	router.HandleFunc("/packages/delete/{packageID:[0-9]+}", h.handleDeletePackage).Methods(("DELETE"))

}

func (h *Handler) handleListActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.castle.ListActivities()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If no activities found, return a message
	if len(activities) == 0 {
		utils.WriteJSON(w, http.StatusOK, "no reviews found")
		return
	}

	// Return the activities as a JSON response
	utils.WriteJSON(w, http.StatusOK, activities)
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
	var payload types.ActivityPayload
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
	// Get the review ID from the URL parameters
	vars := mux.Vars(r)
	str, ok := vars["activityID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing activity ID"))
		return
	}

	// Convert review ID from string to int
	activityID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid activity ID"))
		return
	}

	// Check if the review exists
	existingActivity, err := h.castle.GetActivityByID(activityID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching activity: %w", err))
		return
	}

	// Attempt to delete the review
	err = h.castle.DeleteActivity(existingActivity.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error deleting activity: %w", err))
		return
	}

	// Return a 200 OK response with a success message
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Activity with ID %d successfully deleted", activityID))
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
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("package with name %s already exists", payload.Name))
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

func (h *Handler) handleGetActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["activityID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing activity ID"))
		return
	}

	activityID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid activity ID"))
		return
	}

	activity, err := h.castle.GetActivityByID(activityID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("activity not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, activity)
}

func (h *Handler) handleUpdateActivity(w http.ResponseWriter, r *http.Request) {
	// Get the activity ID from the URL parameters
	vars := mux.Vars(r)
	str, ok := vars["activityID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing activity ID"))
		return
	}

	// Convert activity ID from string to int
	activityID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid activity ID"))
		return
	}

	// Get JSON payload
	var payload types.ActivityPayload
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

	// Check if the review exists
	existingReview, err := h.castle.GetActivityByID(activityID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("activity not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Update the review
	updateActivity := types.Activity{
		ID:          existingReview.ID,
		Name:        payload.Name,
		Description: payload.Description,
		BasePrice:   payload.BasePrice,
		Hidden:      payload.Hidden,
		Category:    payload.Category,
		FkPackageID: payload.FkPackageID,
	}

	err = h.castle.UpdateActivity(updateActivity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updateActivity)
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
