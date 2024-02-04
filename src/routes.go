package routes

import (
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Добро пожаловать на главную страницу!"))
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет, мир!"))
	})

	r.Post("/createuser", func(w http.ResponseWriter, r *http.Request) {
		createUserHandler(w, r, db)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		log.Printf("User %s not found", newUser.Username)
	}

	err = models.CreateUser(db, &newUser)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
