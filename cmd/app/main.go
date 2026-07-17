package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application/usecases/authenticate"
	"github.com/EugeneNail/lifeline/internal/application/usecases/create_completable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/create_entry"
	"github.com/EugeneNail/lifeline/internal/application/usecases/create_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/create_time_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/get_completable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/get_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/get_time_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/list_habits"
	"github.com/EugeneNail/lifeline/internal/application/usecases/refresh"
	"github.com/EugeneNail/lifeline/internal/application/usecases/register_user"
	"github.com/EugeneNail/lifeline/internal/application/usecases/update_completable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/update_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/update_time_habit"
	"github.com/EugeneNail/lifeline/internal/domain/entries"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	"github.com/EugeneNail/lifeline/internal/infrastructure/encryption"
	"github.com/EugeneNail/lifeline/internal/infrastructure/postgres"
	"github.com/EugeneNail/lifeline/internal/infrastructure/tokens"
	transportAuthenticate "github.com/EugeneNail/lifeline/internal/presentation/http/api/authenticate"
	transportCreate_completable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/create_completable_habit"
	transportCreate_entry "github.com/EugeneNail/lifeline/internal/presentation/http/api/create_entry"
	transportCreate_measurable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/create_measurable_habit"
	transportCreate_time_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/create_time_habit"
	transportGet_completable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/get_completable_habit"
	transportGet_measurable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/get_measurable_habit"
	transportGet_time_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/get_time_habit"
	transportList_habits "github.com/EugeneNail/lifeline/internal/presentation/http/api/list_habits"
	transportRefresh "github.com/EugeneNail/lifeline/internal/presentation/http/api/refresh"
	transportRegister_user "github.com/EugeneNail/lifeline/internal/presentation/http/api/register_user"
	transportUpdate_completable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/update_completable_habit"
	transportUpdate_measurable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/update_measurable_habit"
	transportUpdate_time_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/update_time_habit"
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
	requestIdentity := authentication.NewRequestIdentity()

	// --- Section: Usecase handlers ---
	accountRepository, err := postgres.NewAccountRepository(db)
	if err != nil {
		log.Fatalf("creating an account repository: %v", err)
	}

	entryRepository, err := postgres.NewEntryRepository(db)
	if err != nil {
		log.Fatalf("creating an entry repository: %v", err)
	}

	completableHabitRepository, err := postgres.NewCompletableHabitRepository(db)
	if err != nil {
		log.Fatalf("creating a completable habit repository: %v", err)
	}

	measurableHabitRepository, err := postgres.NewMeasurableHabitRepository(db)
	if err != nil {
		log.Fatalf("creating a measurable habit repository: %v", err)
	}

	timeHabitRepository, err := postgres.NewTimeHabitRepository(db)
	if err != nil {
		log.Fatalf("creating a time habit repository: %v", err)
	}

	entryCreationPolicy := entries.NewEntryCreationPolicy(entryRepository)
	habitCreationPolicy := habits.NewHabitCreationPolicy(completableHabitRepository, measurableHabitRepository, timeHabitRepository)
	habitModificationPolicy := habits.NewModificationPolicy()

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

	createEntryUsecase, err := create_entry.NewHandler(entryRepository, entryCreationPolicy)
	if err != nil {
		log.Fatalf("creating a create-entry usecase: %v", err)
	}

	createCompletableHabitUsecase, err := create_completable_habit.NewHandler(completableHabitRepository, habitCreationPolicy)
	if err != nil {
		log.Fatalf("creating a create-completable-habit usecase: %v", err)
	}

	createMeasurableHabitUsecase, err := create_measurable_habit.NewHandler(measurableHabitRepository, habitCreationPolicy)
	if err != nil {
		log.Fatalf("creating a create-measurable-habit usecase: %v", err)
	}

	createTimeHabitUsecase, err := create_time_habit.NewHandler(timeHabitRepository, habitCreationPolicy)
	if err != nil {
		log.Fatalf("creating a create-time-habit usecase: %v", err)
	}

	getCompletableHabitUsecase, err := get_completable_habit.NewHandler(completableHabitRepository)
	if err != nil {
		log.Fatalf("creating a get-completable-habit usecase: %v", err)
	}

	getTimeHabitUsecase, err := get_time_habit.NewHandler(timeHabitRepository)
	if err != nil {
		log.Fatalf("creating a get-time-habit usecase: %v", err)
	}

	getMeasurableHabitUsecase, err := get_measurable_habit.NewHandler(measurableHabitRepository)
	if err != nil {
		log.Fatalf("creating a get-measurable-habit usecase: %v", err)
	}

	listHabitsUsecase, err := list_habits.NewHandler(completableHabitRepository, timeHabitRepository, measurableHabitRepository)
	if err != nil {
		log.Fatalf("creating a list-habits usecase: %v", err)
	}

	updateCompletableHabitUsecase, err := update_completable_habit.NewHandler(completableHabitRepository, habitModificationPolicy)
	if err != nil {
		log.Fatalf("creating an update-completable-habit usecase: %v", err)
	}

	updateTimeHabitUsecase, err := update_time_habit.NewHandler(timeHabitRepository, habitModificationPolicy)
	if err != nil {
		log.Fatalf("creating an update-time-habit usecase: %v", err)
	}

	updateMeasurableHabitUsecase, err := update_measurable_habit.NewHandler(measurableHabitRepository, habitModificationPolicy)
	if err != nil {
		log.Fatalf("creating an update-measurable-habit usecase: %v", err)
	}

	// --- Section: HTTP endpoint handlers ---
	registerUserEndpoint := transportRegister_user.NewHandler(registerUserUsecase)
	authenticateEndpoint := transportAuthenticate.NewHandler(authenticateUsecase)
	refreshEndpoint := transportRefresh.NewHandler(refreshUsecase)
	createEntryEndpoint := transportCreate_entry.NewHandler(createEntryUsecase, requestIdentity)
	createCompletableHabitEndpoint := transportCreate_completable_habit.NewHandler(createCompletableHabitUsecase, requestIdentity)
	createMeasurableHabitEndpoint := transportCreate_measurable_habit.NewHandler(createMeasurableHabitUsecase, requestIdentity)
	listHabitsEndpoint := transportList_habits.NewHandler(listHabitsUsecase, requestIdentity)
	createTimeHabitEndpoint := transportCreate_time_habit.NewHandler(createTimeHabitUsecase, requestIdentity)
	getCompletableHabitEndpoint := transportGet_completable_habit.NewHandler(getCompletableHabitUsecase, requestIdentity)
	getTimeHabitEndpoint := transportGet_time_habit.NewHandler(getTimeHabitUsecase, requestIdentity)
	getMeasurableHabitEndpoint := transportGet_measurable_habit.NewHandler(getMeasurableHabitUsecase, requestIdentity)
	updateCompletableHabitEndpoint := transportUpdate_completable_habit.NewHandler(updateCompletableHabitUsecase, requestIdentity)
	updateMeasurableHabitEndpoint := transportUpdate_measurable_habit.NewHandler(updateMeasurableHabitUsecase, requestIdentity)
	updateTimeHabitEndpoint := transportUpdate_time_habit.NewHandler(updateTimeHabitUsecase, requestIdentity)

	// --- Section: HTTP server ---
	server := http.NewServeMux()
	server.Handle("POST /api/v1/users/register", middleware.WriteJSONResponse(registerUserEndpoint))
	server.Handle("POST /api/v1/users/login", middleware.WriteJSONResponse(authenticateEndpoint))
	server.Handle("POST /api/v1/users/refresh", middleware.WriteJSONResponse(refreshEndpoint))
	server.Handle("POST /api/v1/entries", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createEntryEndpoint)))
	server.Handle("POST /api/v1/habits/completable", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createCompletableHabitEndpoint)))
	server.Handle("POST /api/v1/habits/measurable", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createMeasurableHabitEndpoint)))
	server.Handle("POST /api/v1/habits/time", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createTimeHabitEndpoint)))
	server.Handle("GET /api/v1/habits/completable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getCompletableHabitEndpoint)))
	server.Handle("GET /api/v1/habits/time/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getTimeHabitEndpoint)))
	server.Handle("GET /api/v1/habits/measurable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getMeasurableHabitEndpoint)))
	server.Handle("GET /api/v1/habits", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(listHabitsEndpoint)))
	server.Handle("PUT /api/v1/habits/completable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(updateCompletableHabitEndpoint)))
	server.Handle("PUT /api/v1/habits/measurable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(updateMeasurableHabitEndpoint)))
	server.Handle("PUT /api/v1/habits/time/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(updateTimeHabitEndpoint)))

	// TODO handle the error
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configuration.App.Port), server)
}
