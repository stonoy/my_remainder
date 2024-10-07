package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stonoy/my_remainder/internal/database"
)

type apiConfig struct {
	dbQ        *database.Queries
	hits       int
	jwt_secret string
}

func main() {
	// load local enviroment variables in our program
	err := godotenv.Load()
	if err != nil {
		log.Printf("cound not load local enviroment variables in our program -> %v", err)
	}

	// get the variables
	port := os.Getenv("PORT")
	if port == "" {
		log.Panic("No port provided!")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Panic("No jwt secret provided")
	}

	// initiate apiConfig
	apiCfg := &apiConfig{
		hits:       0,
		jwt_secret: jwtSecret,
	}

	db_uri := os.Getenv("DB_URI")
	if db_uri != "" {
		dbtxType, err := sql.Open("postgres", db_uri)
		if err != nil {
			log.Panicf("Error in connecting to server -> %v", err)
		}

		db_queries := database.New(dbtxType)
		apiCfg.dbQ = db_queries
	}

	// set new router
	mainRouter := chi.NewRouter()

	// basic cors
	mainRouter.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	apiRouter := chi.NewRouter()

	// user
	apiRouter.Post("/register", apiCfg.registerUser)
	apiRouter.Post("/login", apiCfg.login)

	// remainder
	apiRouter.Post("/createremainders", apiCfg.authTokenToUser(apiCfg.createRemainders))
	apiRouter.Get("/getremainders", apiCfg.authTokenToUser(apiCfg.getRemaindersByUser))
	apiRouter.Get("/getremainder/{ID}", apiCfg.authTokenToUser(apiCfg.getRemainderByID))
	apiRouter.Put("/updateremainder/{ID}", apiCfg.authTokenToUser(apiCfg.updateRemainder))
	apiRouter.Delete("/deleteremainder/{ID}", apiCfg.authTokenToUser(apiCfg.deleteRemainder))

	mainRouter.Mount("/api/v1", apiRouter)

	myServer := &http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	log.Printf("Server listening on port %v", port)

	log.Fatal(myServer.ListenAndServe())
}
