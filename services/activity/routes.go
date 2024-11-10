package activity

// TODO: Use sqlx
import (
	"database/sql"
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
	activityCastle types.ActivityCastle
	userCastle     types.UserCastle
}

func NewHandler(activityCastle types.ActivityCastle, userCastle types.UserCastle) *Handler {
	return &Handler{
		activityCastle: activityCastle,
		userCastle:     userCastle}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/activities", h.handleListActivities).Methods(("GET"))
	router.HandleFunc("/activities/create", auth.WithJWTAuth(h.handleCreateActivity, h.userCastle, "administrator", "organizer")).Methods(("POST"))
	router.HandleFunc("/activities/{activityID:[0-9]+}", h.handleGetActivity).Methods(("GET"))
	router.HandleFunc("/activities/update/{activityID:[0-9]+}", auth.WithJWTAuth(h.handleUpdateActivity, h.userCastle, "administrator", "organizer")).Methods(("PUT"))
	router.HandleFunc("/activities/delete/{activityID:[0-9]+}", auth.WithJWTAuth(h.handleDeleteActivity, h.userCastle, "administrator", "organizer")).Methods(("DELETE"))

	router.HandleFunc("/activities/filter", h.handleFilterActivities).Methods(("GET"))

	router.HandleFunc("/packages/create", auth.WithJWTAuth(h.handleCreatePackage, h.userCastle, "administrator", "organizer")).Methods(("POST"))
	router.HandleFunc("/packages/delete/{packageID:[0-9]+}", auth.WithJWTAuth(h.handleDeletePackage, h.userCastle, "administrator", "organizer")).Methods(("DELETE"))

}

// ListActivities godoc
// @Summary      List all activities
// @Description  Returns list of all registered activities
// @Tags         activity
// @Produce      json
// @Success      200  {array}    types.Activity
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /activities [get]
func (h *Handler) handleListActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.activityCastle.ListActivities()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If no activities found, return an empty array
	if len(activities) == 0 {
		utils.WriteJSON(w, http.StatusOK, []types.Activity{})
		return
	}

	// Return the activities as a JSON response
	utils.WriteJSON(w, http.StatusOK, activities)
}

// CreateActivity godoc
// @Summary      Create a new activity
// @Description  Create a new activity with the given name, description, category, price, and package ID
// @Tags         activity
// @Produce      json
// @Param        payload body types.ActivityPayload true "Activity data"
// @Success      201  {object}   types.ErrorResponse "Activity %s successfully created"
// @Failure      400  {object}   types.ErrorResponse "Invalid payload"
// @Failure      422  {object}   types.ErrorResponse "activity with name %s inside package already exists"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /activities/create [post]
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
	_, err := h.activityCastle.GetActivityInsidePackageByName(payload.Name, payload.FkPackageID)
	if err == nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("activity with name %s inside package already exists", payload.Name))
		return
	}

	err = h.activityCastle.CreateActivity(types.Activity{
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

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Activity %s successfully created", payload.Name))
}

// GetActivity godoc
// @Summary      Get activity by ID
// @Description  Get activity data by ID from the database
// @Tags         activity
// @Produce      json
// @Param        activityID path int true "Activity ID"
// @Success      200  {object}   types.Activity
// @Failure      400  {object}   types.ErrorResponse "missing or invalid activity ID"
// @NotFound     404  {object}   types.ErrorResponse "Activity not found"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /activities/{activityID} [get]
func (h *Handler) handleGetActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["activityID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing or invalid activity ID"))
		return
	}

	activityID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing or invalid activity ID"))
		return
	}

	activity, err := h.activityCastle.GetActivityByID(activityID)
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

// UpdateActivity godoc
// @Summary      Update activity by ID
// @Description  Update activity data by ID and specifying the new values
// @Tags         activity
// @Produce      json
// @Param        activityID path int true "Activity ID"
// @Param        payload body types.ActivityPayload true "Activity data"
// @Success      200  {object}   types.Activity
// @Failure      400  {object}   types.ErrorResponse "missing or invalid activity ID"
// @Failure      404  {object}   types.ErrorResponse "Activity not found"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /activities/update/{activityID} [put]
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
	existingReview, err := h.activityCastle.GetActivityByID(activityID)
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

	err = h.activityCastle.UpdateActivity(updateActivity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updateActivity)
}

// DeleteActivity godoc
// @Summary      Delete activity by ID
// @Description  Delete activity data by ID from database
// @Tags         activity
// @Produce      json
// @Param        activityID path int true "Activity ID"
// @Success      200  {object}   types.ErrorResponse "Activity with ID %d successfully deleted"
// @Failure      400  {object}   types.ErrorResponse "missing or invalid activity ID"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /activities/delete/{activityID} [delete]
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
	existingActivity, err := h.activityCastle.GetActivityByID(activityID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Activity with ID %d successfully deleted", activityID))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching activity: %w", err))
		return
	}

	// Attempt to delete the review
	err = h.activityCastle.DeleteActivity(existingActivity.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error deleting activity: %w", err))
		return
	}

	// Return a 200 OK response with a success message
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Activity with ID %d successfully deleted", activityID))
}

// FilterActivities godoc
// @Summary      Filter activities
// @Description  Filter activities by category, rating, price, and hidden status
// @Tags         activity
// @Produce      json
// @Param        payload body types.ActivityFilterPayload true "Filter payload"
// @Success      200  {array}    types.Activity
// @Failure      400  {object}   types.ErrorResponse "Invalid payload"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /activities/filter [get]
func (h *Handler) handleFilterActivities(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.ActivityFilterPayload
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

	activities, err := h.activityCastle.FilterActivities(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If no activities found, return a message
	if len(activities) == 0 {
		utils.WriteJSON(w, http.StatusOK, []types.Activity{})
		return
	}

	// Return the activities as a JSON response
	utils.WriteJSON(w, http.StatusOK, activities)
}

// CreatePackage godoc
// @Summary      Create a new package
// @Description  Create a new package with the given name, description, price, and organizer ID
// @Tags         package
// @Produce      json
// @Param        payload body types.CreatePackagePayload true "Package data"
// @Success      201  {object}   types.ErrorResponse "Package %s successfully created"
// @Failure      400  {object}   types.ErrorResponse "Invalid payload"
// @Failure      422  {object}   types.ErrorResponse "Package with name %s already exists"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /packages/create [post]
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
	_, err := h.activityCastle.GetPackageByName(payload.Name)
	if err == nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("package with name %s already exists", payload.Name))
		return
	}

	// if not create
	err = h.activityCastle.CreatePackage(types.Package{
		Name:          payload.Name,
		Description:   payload.Description,
		Price:         payload.Price,
		FkOrganizerID: payload.FkOrganizerID,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Package  %s successfully created", payload.Name))
}

// DeletePackage godoc
// @Summary      Delete package by ID
// @Description  Delete package data with all activties by ID from database
// @Tags         package
// @Produce      json
// @Param        packageID path int true "Package ID"
// @Success      200  {object}   types.ErrorResponse "Package with ID %d successfully deleted"
// @BadRequest   400  {object}   types.ErrorResponse "missing or invalid package ID"
// @Failure      500  {object}   types.ErrorResponse "Internal server error"
// @Router       /packages/delete/{packageID} [delete]
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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing or invalid package ID %v", errors))
		return
	}

	// delete package
	err := h.activityCastle.DeletePackage(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete package: %v", err))
		return
	}

	// TODO: status NO content
	utils.WriteJSON(w, http.StatusAccepted, fmt.Sprintf("package with id %d successfully deleted", payload.ID))
}
