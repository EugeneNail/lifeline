package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application/usecases/authenticate"
	"github.com/EugeneNail/lifeline/internal/application/usecases/refresh"
	"github.com/EugeneNail/lifeline/internal/application/usecases/register_user"
	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	"github.com/EugeneNail/lifeline/internal/infrastructure/encryption"
	"github.com/EugeneNail/lifeline/internal/infrastructure/postgres"
	"github.com/EugeneNail/lifeline/internal/infrastructure/tokens"
	transportAuthenticate "github.com/EugeneNail/lifeline/internal/presentation/http/api/authenticate"
	transportRefresh "github.com/EugeneNail/lifeline/internal/presentation/http/api/refresh"
	transportRegister_user "github.com/EugeneNail/lifeline/internal/presentation/http/api/register_user"
	"github.com/EugeneNail/lifeline/internal/presentation/http/middleware"
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

	jwtProvider, err := tokens.NewJWTProvider(configuration.JWT.Secret)
	if err != nil {
		log.Fatalf("creating a JWT provider: %v", err)
	}

	bcryptPasswordHasher := encryption.NewBcryptPasswordHasher()
	bcryptPasswordVerifier := encryption.NewBcryptPasswordVerifier()

	// --- Section: Usecase handlers ---
	accountRepository, err := postgres.NewAccountRepository(db)
	if err != nil {
		log.Fatalf("creating an account repository: %v", err)
	}

	registerUserUsecase, err := register_user.NewHandler(bcryptPasswordHasher, accountRepository)
	if err != nil {
		log.Fatalf("creating a register-user usecase: %v", err)
	}

	authenticateUsecase, err := authenticate.NewHandler(accountRepository, bcryptPasswordVerifier, jwtProvider, configuration.App.Environment)
	if err != nil {
		log.Fatalf("creating an authenticate usecase: %v", err)
	}

	refreshUsecase, err := refresh.NewHandler(accountRepository, jwtProvider)
	if err != nil {
		log.Fatalf("creating a refresh usecase: %v", err)
	}

	// --- Section: HTTP endpoint handlers ---
	registerUserEndpoint := transportRegister_user.NewHandler(registerUserUsecase)
	authenticateEndpoint := transportAuthenticate.NewHandler(authenticateUsecase)
	refreshEndpoint := transportRefresh.NewHandler(refreshUsecase)

	// --- Section: HTTP server ---
	server := http.NewServeMux()
	server.Handle("POST /api/v1/users/register", middleware.WriteJSONResponse(registerUserEndpoint))
	server.Handle("POST /api/v1/users/login", middleware.WriteJSONResponse(authenticateEndpoint))
	server.Handle("POST /api/v1/users/refresh", middleware.WriteJSONResponse(refreshEndpoint))
	// TODO handle the error
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configuration.App.Port), server)
}
