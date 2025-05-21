package main

// referência usada para conexão com o BD: https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/

import (
	http "07-twitter/internal/adapters/http"
	handlers "07-twitter/internal/adapters/http/handlers"
	repos "07-twitter/internal/adapters/repositories"
	services "07-twitter/internal/services"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	host, port, user, password, dbname := carregarVariaveisAmbiente()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("[sql.Open()] - Erro ao abrir conexão com o banco de dados: %v \n", err)
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("[db.Ping()] - Erro ao verificar conexão com o banco de dados: %v \n", err)
		panic(err)
	}

	fmt.Println("Sucesso ao se conectar com o Banco de Dados!")

	userRepo := repos.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	router := http.NewRouter(userHandler)
	router.Run(":8080")
}

func carregarVariaveisAmbiente() (host, port, user, password, dbname string) {
	_ = godotenv.Load("../../infra/.env")

	host = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("POSTGRES_PORT")
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname = os.Getenv("POSTGRES_DB")

	return
}
