package main

import (
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application/usecases/auth/authenticate"
	"github.com/EugeneNail/lifeline/internal/application/usecases/auth/refresh"
	"github.com/EugeneNail/lifeline/internal/application/usecases/auth/register_user"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/create_completable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/create_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/create_time_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/get_completable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/get_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/get_time_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/list_habit_records"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/list_habits"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/save_completable_habit_record"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/save_measurable_habit_record"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/save_time_habit_record"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/update_completable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/update_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/update_time_habit"
	"github.com/EugeneNail/lifeline/internal/application/usecases/journal/create_journal"
	"github.com/EugeneNail/lifeline/internal/application/usecases/journal/get_journal"
	"github.com/EugeneNail/lifeline/internal/application/usecases/moods/get_mood_record"
	"github.com/EugeneNail/lifeline/internal/application/usecases/moods/save_mood_record"
	"github.com/EugeneNail/lifeline/internal/application/usecases/transactions/create_transaction"
	transportAuthenticate "github.com/EugeneNail/lifeline/internal/presentation/http/api/auth/authenticate"
	transportRefresh "github.com/EugeneNail/lifeline/internal/presentation/http/api/auth/refresh"
	transportRegister_user "github.com/EugeneNail/lifeline/internal/presentation/http/api/auth/register_user"
	transportCreate_completable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/create_completable_habit"
	transportCreate_measurable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/create_measurable_habit"
	transportCreate_time_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/create_time_habit"
	transportGet_completable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/get_completable_habit"
	transportGet_measurable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/get_measurable_habit"
	transportGet_time_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/get_time_habit"
	transportList_habit_records "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/list_habit_records"
	transportList_habits "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/list_habits"
	transportSave_completable_habit_record "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/save_completable_habit_record"
	transportSave_measurable_habit_record "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/save_measurable_habit_record"
	transportSave_time_habit_record "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/save_time_habit_record"
	transportUpdate_completable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/update_completable_habit"
	transportUpdate_measurable_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/update_measurable_habit"
	transportUpdate_time_habit "github.com/EugeneNail/lifeline/internal/presentation/http/api/habit/update_time_habit"
	transportGet_journal "github.com/EugeneNail/lifeline/internal/presentation/http/api/journal/get_journal"
	transportCreate_journal "github.com/EugeneNail/lifeline/internal/presentation/http/api/journal/save_journal"
	transportGet_mood_record "github.com/EugeneNail/lifeline/internal/presentation/http/api/moods/get_mood_record"
	transportSave_mood_record "github.com/EugeneNail/lifeline/internal/presentation/http/api/moods/save_mood_record"
	transportCreate_transaction "github.com/EugeneNail/lifeline/internal/presentation/http/api/transactions/create_transaction"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	habitrecords "github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
	"github.com/EugeneNail/lifeline/internal/infrastructure/encryption"
	"github.com/EugeneNail/lifeline/internal/infrastructure/postgres"
	"github.com/EugeneNail/lifeline/internal/infrastructure/tokens"
	"github.com/EugeneNail/lifeline/internal/presentation/http/middleware"
)

const frontendAssetsDir = "internal/presentation/http/web/dist"

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

	journalRepository, err := postgres.NewJournalRepository(db)
	if err != nil {
		log.Fatalf("creating a journal repository: %v", err)
	}

	moodRecordRepository, err := postgres.NewRecordRepository(db)
	if err != nil {
		log.Fatalf("creating a mood record repository: %v", err)
	}

	transactionRepository, err := postgres.NewTransactionRepository(db)
	if err != nil {
		log.Fatalf("creating a transaction repository: %v", err)
	}

	completableHabitRepository, err := postgres.NewCompletableHabitRepository(db)
	if err != nil {
		log.Fatalf("creating a completable habit repository: %v", err)
	}

	completableHabitRecordRepository, err := postgres.NewCompletableHabitRecordRepository(db)
	if err != nil {
		log.Fatalf("creating a completable habit record repository: %v", err)
	}

	timeHabitRecordRepository, err := postgres.NewTimeHabitRecordRepository(db)
	if err != nil {
		log.Fatalf("creating a time habit record repository: %v", err)
	}

	measurableHabitRepository, err := postgres.NewMeasurableHabitRepository(db)
	if err != nil {
		log.Fatalf("creating a measurable habit repository: %v", err)
	}

	measurableHabitRecordRepository, err := postgres.NewMeasurableHabitRecordRepository(db)
	if err != nil {
		log.Fatalf("creating a measurable habit record repository: %v", err)
	}

	timeHabitRepository, err := postgres.NewTimeHabitRepository(db)
	if err != nil {
		log.Fatalf("creating a time habit repository: %v", err)
	}

	habitCreationPolicy := habits.NewHabitCreationPolicy(completableHabitRepository, measurableHabitRepository, timeHabitRepository)
	habitModificationPolicy := habits.NewModificationPolicy()
	habitSavingPolicy := habitrecords.NewSavingPolicy()

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

	createJournalUsecase, err := create_journal.NewHandler(journalRepository)
	if err != nil {
		log.Fatalf("creating a create-journal usecase: %v", err)
	}

	getJournalUsecase, err := get_journal.NewHandler(journalRepository)
	if err != nil {
		log.Fatalf("creating a get-journal usecase: %v", err)
	}

	saveMoodRecordUsecase, err := save_mood_record.NewHandler(moodRecordRepository)
	if err != nil {
		log.Fatalf("creating a save-mood-record usecase: %v", err)
	}

	createTransactionUsecase, err := create_transaction.NewHandler(transactionRepository)
	if err != nil {
		log.Fatalf("creating a create-transaction usecase: %v", err)
	}

	getMoodRecordUsecase, err := get_mood_record.NewHandler(moodRecordRepository)
	if err != nil {
		log.Fatalf("creating a get-mood-record usecase: %v", err)
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

	listHabitRecordsUsecase, err := list_habit_records.NewHandler(completableHabitRecordRepository, timeHabitRecordRepository, measurableHabitRecordRepository)
	if err != nil {
		log.Fatalf("creating a list-habit-records usecase: %v", err)
	}

	saveCompletableHabitRecordUsecase, err := save_completable_habit_record.NewHandler(completableHabitRecordRepository, completableHabitRepository, habitSavingPolicy)
	if err != nil {
		log.Fatalf("creating a save-completable-habit-record usecase: %v", err)
	}

	saveTimeHabitRecordUsecase, err := save_time_habit_record.NewHandler(timeHabitRecordRepository, timeHabitRepository, habitSavingPolicy)
	if err != nil {
		log.Fatalf("creating a save-time-habit-record usecase: %v", err)
	}

	saveMeasurableHabitRecordUsecase, err := save_measurable_habit_record.NewHandler(measurableHabitRecordRepository, measurableHabitRepository, habitSavingPolicy)
	if err != nil {
		log.Fatalf("creating a save-measurable-habit-record usecase: %v", err)
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
	createJournalEndpoint := transportCreate_journal.NewHandler(createJournalUsecase, requestIdentity)
	getJournalEndpoint := transportGet_journal.NewHandler(getJournalUsecase, requestIdentity)
	getMoodRecordEndpoint := transportGet_mood_record.NewHandler(getMoodRecordUsecase, requestIdentity)
	saveMoodRecordEndpoint := transportSave_mood_record.NewHandler(saveMoodRecordUsecase, requestIdentity)
	createTransactionEndpoint := transportCreate_transaction.NewHandler(createTransactionUsecase, requestIdentity)
	createCompletableHabitEndpoint := transportCreate_completable_habit.NewHandler(createCompletableHabitUsecase, requestIdentity)
	createMeasurableHabitEndpoint := transportCreate_measurable_habit.NewHandler(createMeasurableHabitUsecase, requestIdentity)
	listHabitsEndpoint := transportList_habits.NewHandler(listHabitsUsecase, requestIdentity)
	listHabitRecordsEndpoint := transportList_habit_records.NewHandler(listHabitRecordsUsecase, requestIdentity)
	saveCompletableHabitRecordEndpoint := transportSave_completable_habit_record.NewHandler(saveCompletableHabitRecordUsecase, requestIdentity)
	saveMeasurableHabitRecordEndpoint := transportSave_measurable_habit_record.NewHandler(saveMeasurableHabitRecordUsecase, requestIdentity)
	saveTimeHabitRecordEndpoint := transportSave_time_habit_record.NewHandler(saveTimeHabitRecordUsecase, requestIdentity)
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
	server.Handle("POST /api/v1/transactions", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createTransactionEndpoint)))
	server.Handle("GET /api/v1/journals/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getJournalEndpoint)))
	server.Handle("POST /api/v1/journals/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createJournalEndpoint)))
	server.Handle("GET /api/v1/moods/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getMoodRecordEndpoint)))
	server.Handle("POST /api/v1/moods/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(saveMoodRecordEndpoint)))
	server.Handle("POST /api/v1/habits/completable", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createCompletableHabitEndpoint)))
	server.Handle("POST /api/v1/habits/measurable", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createMeasurableHabitEndpoint)))
	server.Handle("POST /api/v1/habits/time", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(createTimeHabitEndpoint)))
	server.Handle("GET /api/v1/habits/completable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getCompletableHabitEndpoint)))
	server.Handle("GET /api/v1/habits/time/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getTimeHabitEndpoint)))
	server.Handle("GET /api/v1/habits/measurable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(getMeasurableHabitEndpoint)))
	server.Handle("GET /api/v1/habits", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(listHabitsEndpoint)))
	server.Handle("GET /api/v1/habits/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(listHabitRecordsEndpoint)))
	server.Handle("POST /api/v1/habits/completable/{uuid}/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(saveCompletableHabitRecordEndpoint)))
	server.Handle("POST /api/v1/habits/measurable/{uuid}/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(saveMeasurableHabitRecordEndpoint)))
	server.Handle("POST /api/v1/habits/time/{uuid}/{date}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(saveTimeHabitRecordEndpoint)))
	server.Handle("PUT /api/v1/habits/completable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(updateCompletableHabitEndpoint)))
	server.Handle("PUT /api/v1/habits/measurable/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(updateMeasurableHabitEndpoint)))
	server.Handle("PUT /api/v1/habits/time/{uuid}", middleware.Authenticate(jwtProvider, requestIdentity)(middleware.WriteJSONResponse(updateTimeHabitEndpoint)))
	server.Handle("/", newPublicRoutesHandler(frontendAssetsDir))

	// TODO handle the error
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configuration.App.Port), server); err != nil {
		log.Fatalf("serving HTTP server: %v", err)
	}
}

// newPublicRoutesHandler returns an HTTP handler that serves the built frontend for public routes.
func newPublicRoutesHandler(assetsDir string) http.Handler {
	indexPath := filepath.Join(assetsDir, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		requestPath := normalizeFrontendRequestPath(r.URL.Path)
		if requestPath == "index.html" {
			http.ServeFile(w, r, indexPath)
			return
		}

		assetPath := filepath.Join(assetsDir, filepath.FromSlash(requestPath))
		if fileExists(assetPath) {
			http.ServeFile(w, r, assetPath)
			return
		}

		if isStaticAssetRequest(requestPath) {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, indexPath)
	})
}

// normalizeFrontendRequestPath converts a request path into a safe relative path within the frontend bundle.
func normalizeFrontendRequestPath(requestPath string) string {
	cleanedPath := path.Clean(strings.TrimPrefix(requestPath, "/"))
	if cleanedPath == "." {
		return "index.html"
	}

	if strings.HasPrefix(cleanedPath, "..") {
		return "index.html"
	}

	return cleanedPath
}

// fileExists reports whether the given file exists on disk.
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// isStaticAssetRequest reports whether the request is looking for a direct static asset.
func isStaticAssetRequest(requestPath string) bool {
	baseName := path.Base(requestPath)
	return strings.Contains(baseName, ".")
}
