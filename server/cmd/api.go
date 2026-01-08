package main

import (
	"log"
	"net/http"
	repo "server/internal/adapters/postgresql/sqlc"
	"server/internal/products"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	dbConn *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}
type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	// chi is used for routing, middleware things
	r := chi.NewRouter()
	r.Use(middleware.RequestID) //used for rate limiting
	r.Use(middleware.RealIP)    //used for rate limiting and analytics and tracing
	r.Use(middleware.Logger)    //
	r.Use(middleware.Recoverer)
	// if request takes more than 60 seconds, timeout then just stop
	r.Use(middleware.Timeout(60 * time.Second))
	// user->handler GET /products->service get products->repo SELECT * FROM products
	productService := products.NewService(repo.New(app.dbConn))
	handler := products.NewHandler(productService)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Get("/products", handler.ListProducts)
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         ":" + app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at addr %s ", app.config.addr)
	return srv.ListenAndServe()
}
