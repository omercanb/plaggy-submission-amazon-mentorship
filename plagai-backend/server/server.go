package server

import (
	// standard lib imports
	"fmt"
	"log"
	"net/http"
	"os"

	// external libs
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// internal packages
	"github.com/plagai/plagai-backend/api/routeHandles"
	"github.com/plagai/plagai-backend/middleware"
	"github.com/plagai/plagai-backend/models/database"
)

func Start() {
	godotenv.Load()
	// Connection string for Neon Postgres
	connectionString := os.Getenv("DATABASE_URL")
	// Connect to the database
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set to Info level for development
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Get the underlying SQL DB object
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB object: %v", err)
	}
	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	fmt.Println("Successfully connected to Neon Postgres database!")
	// Auto-migrate the schema
	if os.Getenv("ENV") != "DEV" {
		err = db.AutoMigrate(&database.Assignment{}, &database.Classroom{}, &database.Diff{}, &database.Flag{}, &database.Instructor{}, &database.Student{}, &database.StudentAssignment{})
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		fmt.Println("Database migrated successfully!")
	}

	// dbscripts.HashPasswords(db)
	// dbscripts.Populate(db)

	h := &routeHandles.Handler{DB: db}

	r := mux.NewRouter()
	// API routes are defined here. Handler functions are defined under routeHandles. When we need to define a new route
	// first we say api.HandleFunc("/api-route", handlerFunction).Methods("GET" | "POST" | ...) ~brtcrt
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/login", h.Login).Methods("POST")
	api.HandleFunc("/auth/magic-request", routeHandles.MagicRequestHandler).Methods("POST")
	api.HandleFunc("/auth/magic-status", routeHandles.MagicStatusHandler).Methods("POST")
	api.HandleFunc("/auth/magic-consume", routeHandles.MagicConsumeHandler).Methods("GET")
	// protected endpoints
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/health", routeHandles.HealthCheck).Methods("GET")
	protected.HandleFunc("/submit", h.SubmitHandler).Methods("POST")
	protected.HandleFunc("/assignments", h.SendAssignments).Methods("GET")

	// Currently giving a JWT token timed out error and will ask brtcrt about it later
	// protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/detections", h.SendDetections).Methods("GET")
	protected.HandleFunc("/sections", h.SendSections).Methods("GET")
	protected.HandleFunc("/homeworks", h.SendHomeworks).Methods("GET")
	protected.HandleFunc("/section", h.SendSectionDetails).Methods("GET")
	protected.HandleFunc("/homework", h.SendHomeworkDetails).Methods("GET")
	protected.HandleFunc("/create_homeworks", h.CreateHomework).Methods("POST")
	protected.HandleFunc("/build_file", h.BuildFile).Methods("GET")
	protected.HandleFunc("/homework/students", h.ListHomeworkStudents).Methods("GET")
	protected.HandleFunc("/homework/files", h.ListStudentFiles).Methods("GET")
	// What is this?
	/*
		In very simple terms, this is a method of disallowing cross origin request forgery. What this should
		hopefully do is prevent access to these api endpoint unless the request is coming from a specific and
		authenticated origin. In our case that origin should be localhost:3000 or the same origin in production
		since idealy, this server should also serve our compiled frontend code. Although that is for much later.
		So for now, this should work just fine. ~brtcrt
		jk. I nuked CSRF from this. For now, CORS should suffice. ~brtcrt
	*/
	// this is only for dev env. Change this in production to "/"? I think that is how it worked. ~brtcrt
	// 2 months later and I still have no clue how this works. I just added these in the server manually
	// and it decided to work so now I am too scared to change it. ~brtcrt
	allowedOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://13.51.70.165:3000", "http://plaggy.xyz", "/"}
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-CSRF-Token"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	headersExposed := handlers.ExposedHeaders([]string{"X-Total-Count", "X-Page", "X-Limit", "X-Has-More", "X-Next-Page", "token"})

	// Also what are middlewares?
	/*
		If I remember correctly, middlewares are function that our server puts every incoming and outgoing request through.
		In our case our csrf middleware validates that the incoming request is in fact coming from an allowed origin and also
		tags any outgoing request with X-CSRF-Token.
	*/
	handler := handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		headersOk,
		methodsOk,
		headersExposed,
		handlers.AllowCredentials(),
	)(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Backend server starting on port %s", port)
	// log.Fatal(http.ListenAndServe(":"+port, handler))

	loggedRouter := handlers.LoggingHandler(os.Stdout, handler)
	// fucking 0.0.0.0 is required? I don't fucking know anymore ~brtcrt
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, loggedRouter))
}
