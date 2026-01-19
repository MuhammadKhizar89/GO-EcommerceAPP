package main

import (
	"log"
	"net/http"
	repo "server/internal/adapters/postgresql/sqlc"
	"server/internal/domain/auth"
	orders "server/internal/domain/order"
	products "server/internal/domain/product"
	"server/internal/transport/https/handlers"
	customMiddleware "server/internal/transport/https/handlers/middleware"
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
	r.Use(customMiddleware.CORSMiddleware) // Enable CORS

	// user->handler GET /products->service get products->repo SELECT * FROM products
	authService := auth.NewAuthService(repo.New(app.dbConn))
	orderService := orders.NewService(repo.New(app.dbConn), app.dbConn)
	productService := products.NewService(repo.New(app.dbConn))
	productHandler := handlers.NewProductHandler(productService)
	authHandler := handlers.NewAuthHandler(authService)
	ordersHandler := handlers.NewOrderHandler(orderService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Post("/login", authHandler.Login)
	r.Post("/signup", authHandler.Signup)
	r.Get("/products", productHandler.ListProducts)
	r.Post("/products", productHandler.CreateProduct)
	r.Route("/orders", func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		r.Post("/", ordersHandler.PlaceOrder)
		r.Get("/", ordersHandler.GetAllOrders)
	})
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
