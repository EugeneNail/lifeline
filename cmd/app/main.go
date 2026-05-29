package main

import (
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application/usecases/register_user"
	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	"github.com/EugeneNail/lifeline/internal/infrastructure/encryption"
	"github.com/EugeneNail/lifeline/internal/infrastructure/postgres"
	transportRegister_user "github.com/EugeneNail/lifeline/internal/presentation/http/api/register_user"
	"github.com/EugeneNail/lifeline/internal/presentation/http/middleware"
	"log"
	"net/http"
)

func main() {
	// --- Section: Configuration ---
	configuration, err := config.Load()
	if err != nil {
		log.Fatalf("loading a configuration instance: %v", err)
	}

	// --- Section: Infrastructure ---
	db, err := postgres.Connect(configuration.Database.Postgres)
	if err != nil {
		log.Fatalf("connecting to the database: %v", err)
	}

	// --- Section: Usecase handlers ---
	bcryptPasswordHasher := encryption.NewBcryptPasswordHasher()
	//TODO handle the errors
	accountRepository, err := postgres.NewAccountRepository(db)
	registerUserUsecase, err := register_user.NewHandler(bcryptPasswordHasher, accountRepository)
	registerUserEndpoint := transportRegister_user.NewHandler(registerUserUsecase)

	// --- Section: HTTP server ---
	server := http.NewServeMux()
	server.Handle("POST /api/v1/users/register", middleware.WriteJSONResponse(registerUserEndpoint))
	// TODO handle the error
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configuration.App.Port), server)
}
