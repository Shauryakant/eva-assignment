package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"ticket-system/database"
	"ticket-system/middleware"
	"ticket-system/models"
	"ticket-system/utils"
)

type AuthHandler struct {
	db        *database.DB
	jwtSecret string
}

func NewAuthHandler(db *database.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	if email == "" || password == "" || !strings.Contains(email, "@") {
		utils.JSONError(w, http.StatusBadRequest, "valid email and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	usersColl := h.db.UsersCollection()

	var existingUser models.User
	err := usersColl.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		utils.JSONError(w, http.StatusConflict, "user with this email already exists")
		return
	} else if err != mongo.ErrNoDocuments {
		utils.JSONError(w, http.StatusInternalServerError, "database error")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to process password")
		return
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now().UTC(),
	}

	_, err = usersColl.InsertOne(ctx, newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			utils.JSONError(w, http.StatusConflict, "user with this email already exists")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"id":      newUser.ID.Hex(),
		"email":   newUser.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	if email == "" || password == "" {
		utils.JSONError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	usersColl := h.db.UsersCollection()

	var user models.User
	err := usersColl.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.JSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		utils.JSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	claims := &middleware.JWTClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}
