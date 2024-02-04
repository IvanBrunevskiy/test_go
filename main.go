package main

import (
	routes "blog-server/src"
	config "blog-server/src/config"
	"log"
	"path/filepath"

	//"blog-server/src/models"
	//"database/sql"
	//"fmt"
	//"log"
	_ "github.com/lib/pq"

	//"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	configPath := filepath.Join("config.json") // Укажите правильный путь к файлу config.json
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config failed: %v", err)
	}
	//
	//psqlInfo := "host=localhost user=postgres dbname=test_ivan sslmode=disable password=password"
	//
	//db, err := sql.Open("postgres", psqlInfo)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(db)
	//defer db.Close()
	//
	//newUser := models.User{
	//	Username: "newuser",
	//	Password: "password",
	//}
	//err = models.CreateUser(db, &newUser)
	//if err != nil {
	//	log.Fatalf("Failed to create user: %v", err)
	//}
	//
	//log.Println("User created successfully")

	r := routes.Routes(cfg) // Используйте функцию Routes для получения маршрутизатора

	// Использование мидлвара для логирования HTTP-запросов
	//r.Use(middleware.Logger)

	// Запуск сервера на порту 8080
	http.ListenAndServe(":8080", r)
}
