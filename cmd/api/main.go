package main

import (
	"log"
	"net/http"

	authadapters "github.com/thamesblanq/healthconnect/internal/auth/adapters"
	jwtadapter "github.com/thamesblanq/healthconnect/internal/auth/adapters/jwt"
	authapplication "github.com/thamesblanq/healthconnect/internal/auth/application"
	authhandlers "github.com/thamesblanq/healthconnect/internal/auth/handlers"

	"github.com/thamesblanq/healthconnect/internal/config"
	"github.com/thamesblanq/healthconnect/internal/database"

	"github.com/thamesblanq/healthconnect/internal/security/adapters/argon2"

	userpostgres "github.com/thamesblanq/healthconnect/internal/user/adapters/postgres"
	userapplication "github.com/thamesblanq/healthconnect/internal/user/application"
	userhandlers "github.com/thamesblanq/healthconnect/internal/user/handlers"
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

	userRepository := userpostgres.NewUserRepository(db)

	passwordHasher := argon2.NewPasswordHasher()

	userProvider := authadapters.NewUserProvider(userRepository)

	tokenGenerator := jwtadapter.NewTokenGenerator(
		cfg.JWTSecret,
		cfg.JWTExpiration,
	)

	tokenVerifier := jwtadapter.NewTokenVerifier(
		cfg.JWTSecret,
	)

	// --------------------------------------------------
	// 4. Create application use cases
	// --------------------------------------------------

	registerUserUseCase := userapplication.NewRegisterUserUseCase(
		userRepository,
		passwordHasher,
	)

	loginUserUseCase := authapplication.NewLoginUserUseCase(
		userProvider,
		passwordHasher,
		tokenGenerator,
	)

	// --------------------------------------------------
	// 5. Create HTTP handlers
	// --------------------------------------------------

	userHandler := userhandlers.NewHandler(registerUserUseCase)

	authHandler := authhandlers.NewHandler(loginUserUseCase)

	authMiddleware := authhandlers.NewAuthMiddleware(
		tokenVerifier,
	)

	// --------------------------------------------------
	// 6. HTTP server
	// --------------------------------------------------

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("HealthConnect API is running"))
	})

	// User routes
	mux.HandleFunc("/users", userHandler.RegisterUser)

	// Auth routes
	mux.HandleFunc("/auth/login", authHandler.Login)

	mux.Handle(
		"/users/me",
		authMiddleware.RequireAuth(
			http.HandlerFunc(userHandler.GetMe),
		),
	)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server running on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
