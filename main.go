package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// User – структура пользователя
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Хранилище пользователей в памяти
var (
	users   = make(map[string]User) // ключ – email
	usersMu sync.Mutex
)

// Секрет для подписи JWT
var jwtKey = []byte("my_secret_key")

// Claims – структура для JWT
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func main() {
	// Инициализируем пользователя admin
	usersMu.Lock()
	users["admin@example.com"] = User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: "admin123",
		Role:     "admin",
	}
	usersMu.Unlock()

	// Обработчики эндпоинтов
	http.Handle("/", http.FileServer(http.Dir("./static"))) // предполагаем, что файлы index.html, style.css и app.js находятся в папке static
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// registerHandler обрабатывает регистрацию
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Если пользователь уже существует, выдаем ошибку
	if _, exists := users[newUser.Email]; exists {
		http.Error(w, "Пользователь уже существует", http.StatusBadRequest)
		return
	}

	// При регистрации устанавливаем роль "student"
	newUser.Role = "student"
	users[newUser.Email] = newUser

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Регистрация прошла успешно"))
}

// loginHandler обрабатывает вход пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	user, exists := users[creds.Email]
	usersMu.Unlock()

	if !exists || user.Password != creds.Password {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// Создаем JWT с данными пользователя, устанавливаем время жизни токена (например, 1 час)
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON с токеном и ролью
	response := map[string]string{
		"token": tokenString,
		"role":  user.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
