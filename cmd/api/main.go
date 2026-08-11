package main

import (
	"log"
	"net/http"

	"github.com/thamesblanq/healthconnect/internal/config"
	"github.com/thamesblanq/healthconnect/internal/database"

	"github.com/thamesblanq/healthconnect/internal/user/adapters/argon2"
	"github.com/thamesblanq/healthconnect/internal/user/adapters/postgres"
	"github.com/thamesblanq/healthconnect/internal/user/application"
	"github.com/thamesblanq/healthconnect/internal/user/handlers"
)

func main() {
	// --------------------------------------------------
	// 1. Load configuration
	// --------------------------------------------------

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// --------------------------------------------------
	// 2. Connect to PostgreSQL
	// --------------------------------------------------

	db, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// --------------------------------------------------
	// 3. Create adapters
	// --------------------------------------------------

	userRepository := postgres.NewUserRepository(db)

	passwordHasher := argon2.NewPasswordHasher()

	// --------------------------------------------------
	// 4. Create application use cases
	// --------------------------------------------------

	registerUserUseCase := application.NewRegisterUserUseCase(
		userRepository,
		passwordHasher,
	)

	// Prevent the compiler from complaining while we
	// haven't created the HTTP handler yet.
	userHandler := handlers.NewHandler(registerUserUseCase)

	// --------------------------------------------------
	// 5. HTTP server
	// --------------------------------------------------

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("HealthConnect API is running"))
	})
	mux.HandleFunc("/users", userHandler.RegisterUser)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server running on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
