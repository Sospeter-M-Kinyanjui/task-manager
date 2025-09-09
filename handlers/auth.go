package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Sospeter-M-Kinyanjui/task-manager/database"
	"github.com/Sospeter-M-Kinyanjui/task-manager/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashed))
	if err != nil {
		http.Error(w, "User exists or DB error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	var dbUser models.User
	err := database.DB.QueryRow(context.Background(),
		"SELECT id, password FROM users WHERE username=$1", user.Username).Scan(&dbUser.ID, &dbUser.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
