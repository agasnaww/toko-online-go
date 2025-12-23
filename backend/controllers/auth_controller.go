package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"toko-online-go/config"
	"toko-online-go/models"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal mengenkripsi password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	query := "INSERT INTO users (nama, email, password) VALUES (?, ?, ?)"
	_, err = config.DB.Exec(query, user.Nama, user.Email, user.Password)

	if err != nil {
		http.Error(w, "Gagal daftar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Register Berhasil! Silakan Login."})
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userInput models.User

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	var userDB models.User
	query := "SELECT id, nama, email, password FROM users WHERE email = ?"
	err := config.DB.QueryRow(query, userInput.Email).Scan(&userDB.ID, &userDB.Nama, &userDB.Email, &userDB.Password)

	if err == sql.ErrNoRows {
		http.Error(w, "Email atau Password salah", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(userInput.Password))
	if err != nil {
		http.Error(w, "Email atau Password salah", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id": userDB.ID,
		"nama":    userDB.Nama,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login Berhasil!",
		"token":   tokenString,
	})
}
