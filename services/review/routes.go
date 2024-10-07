package review

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
	castle types.ReviewCastle
}

func NewHandler(castle types.ReviewCastle) *Handler {
	return &Handler{castle: castle}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/reviews", h.handleListReviews).Methods(("GET"))
	router.HandleFunc("/reviews/create", h.handleCreateReview).Methods(("POST"))
	router.HandleFunc("/reviews/{reviewID}", h.handleGetReview).Methods(("GET"))
	router.HandleFunc("/reviews/update/{reviewID:[0-9]+}", h.handleUpdateReview).Methods("PUT")
	router.HandleFunc("/reviews/delete/{reviewID:[0-9]+}", h.handleDeleteReview).Methods("DELETE")
}

func (h *Handler) handleListReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.castle.ListReviews()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// If no reviews found, return a message
	if len(reviews) == 0 {
		utils.WriteJSON(w, http.StatusOK, "no reviews found")
		return
	}

	// Return the reviews as a JSON response
	utils.WriteJSON(w, http.StatusOK, reviews)
}

func (h *Handler) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.CreateReviewPayload
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

	// check if the review exists
	_, err := h.castle.GetReviewFromActivityByID(payload.FkActivityID, payload.FkUserID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("review from same user: %s already exists", payload.Comment))
		return
	}

	// if it doesnt create the new review
	err = h.castle.CreateReview(types.Review{
		Comment:      &payload.Comment,
		Rating:       payload.Rating,
		FkUserID:     payload.FkUserID,
		FkActivityID: payload.FkActivityID,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["reviewID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing review ID"))
		return
	}

	reviewID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid review ID"))
		return
	}

	review, err := h.castle.GetReviewByID(reviewID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, review)
}

func (h *Handler) handleUpdateReview(w http.ResponseWriter, r *http.Request) {
	// Get the review ID from the URL parameters
	vars := mux.Vars(r)
	reviewIDStr, ok := vars["reviewID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing review ID"))
		return
	}

	// Convert review ID from string to int
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid review ID"))
		return
	}

	// Get JSON payload
	var payload types.UpdateReviewPayload
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
	existingReview, err := h.castle.GetReviewByID(reviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("review not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Update the review
	updatedReview := types.Review{
		ID:           reviewID,
		Comment:      &payload.Comment,
		Rating:       payload.Rating,
		FkUserID:     existingReview.FkUserID,
		FkActivityID: existingReview.FkActivityID,
	}

	err = h.castle.UpdateReview(updatedReview)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedReview)
}

func (h *Handler) handleDeleteReview(w http.ResponseWriter, r *http.Request) {
	// Get the review ID from the URL parameters
	vars := mux.Vars(r)
	reviewIDStr, ok := vars["reviewID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing review ID"))
		return
	}

	// Convert review ID from string to int
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid review ID"))
		return
	}

	// Check if the review exists
	existingReview, err := h.castle.GetReviewByID(reviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching review: %w", err))
		return
	}

	// Attempt to delete the review
	err = h.castle.DeleteReviewByID(existingReview.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error deleting review: %w", err))
		return
	}

	// Return a 200 OK response with a success message
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Review with ID %d successfully deleted", reviewID))
}
