package main

import (
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application/usecases/register_user"
	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	"github.com/EugeneNail/lifeline/internal/infrastructure/encryption"
	transportRegister_user "github.com/EugeneNail/lifeline/internal/presentation/http/api/register_user"
	"github.com/EugeneNail/lifeline/internal/presentation/http/middleware"
	"log"
	"net/http"
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		log.Fatalf("loading a configuration instance: %v", err)
	}

	bcryptPasswordHasher := encryption.NewBcryptPasswordHasher()
	//TODO handle the error
	registerUserUsecase, err := register_user.NewHandler(bcryptPasswordHasher)
	registerUserEndpoint := transportRegister_user.NewHandler(registerUserUsecase)

	server := http.NewServeMux()
	server.Handle("POST /api/v1/users/register", middleware.WriteJSONResponse(registerUserEndpoint))
	// TODO handle the error
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configuration.App.Port), server)
}
