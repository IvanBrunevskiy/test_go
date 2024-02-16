package routes

import (
	"blog-server/src/auth"
	"blog-server/src/config"
	"blog-server/src/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func Routes(cfg config.Config) *chi.Mux {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName))
	//psqlInfo := "host=localhost user=postgres dbname=test_ivan sslmode=disable password=password"
	//db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)
	//defer db.Close()

	r := chi.NewRouter()

	//r.Use(middleware.AuthMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Добро пожаловать на главную страницу!"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет, мир!"))
	})

	r.Post("/createuser", func(w http.ResponseWriter, r *http.Request) {
		createUserHandler(w, r, db)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, db)
	})

	return r
}

func createUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var newUser models.User
	err1 := json.NewDecoder(r.Body).Decode(&newUser)
	if err1 != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByUsername(db, newUser.Username)
	if err != nil {
		log.Fatalf("Error getting user11: %v", err)
	}
	if user != nil {
		http.Error(w, "User already exist", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = models.CreateUser(db, &newUser)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func loginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var credentials models.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByUsername(db, credentials.Username)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil || user.Password != credentials.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Логин и пароль верны, генерируем токен
	token, err := auth.GenerateToken(credentials.Username)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Сохраняем токен в базе данных (предположим, у вас есть соответствующее поле в вашей таблице)
	//err = models.UpdateUserToken(db, credentials.Username, token)
	//if err != nil {
	//	log.Fatalf("Error updating user token: %v", err)
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	return
	//}

	// Отправляем токен в ответе
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
