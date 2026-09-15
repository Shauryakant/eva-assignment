package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"ticket-system/database"
	"ticket-system/middleware"
	"ticket-system/models"
	"ticket-system/utils"
)

type TicketHandler struct {
	db *database.DB
}

func NewTicketHandler(db *database.DB) *TicketHandler {
	return &TicketHandler{
		db: db,
	}
}

func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input models.CreateTicketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		utils.JSONError(w, http.StatusBadRequest, "title is required")
		return
	}

	now := time.Now().UTC()
	ticket := models.Ticket{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		Status:      models.StatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.db.TicketsCollection().InsertOne(ctx, ticket)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create ticket")
		return
	}

	utils.JSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cursor, err := h.db.TicketsCollection().Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to fetch tickets")
		return
	}
	defer cursor.Close(ctx)

	tickets := make([]models.Ticket, 0)
	if err := cursor.All(ctx, &tickets); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to decode tickets")
		return
	}

	utils.JSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	ticketID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, "ticket not found")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var ticket models.Ticket
	err = h.db.TicketsCollection().FindOne(ctx, bson.M{"_id": ticketID, "user_id": userID}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.JSONError(w, http.StatusNotFound, "ticket not found")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "database error")
		return
	}

	utils.JSON(w, http.StatusOK, ticket)
}
