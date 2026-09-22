package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SephirothGit/warehouse/internal/auth"
	"github.com/SephirothGit/warehouse/internal/handler"
	"github.com/SephirothGit/warehouse/internal/repository"
	"github.com/SephirothGit/warehouse/internal/service"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("unable to open db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("unable to connect to db", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database successfully")

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	userRepo := repository.NewUserRepo(db)
	refreshRepo := repository.NewRefreshTokenRepo(db)
	warehouseRepo := repository.NewWarehouseRepo(db)
	zoneRepo := repository.NewZoneRepo(db)
	rackRepo := repository.NewRackRepo(db)
	shelfRepo := repository.NewShelfRepo(db)
	productRepo := repository.NewProductRepo(db)
	stockRepo := repository.NewStockRepo(db)

	authService := service.NewAuthService(userRepo, refreshRepo, jwtSecret)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	zoneService := service.NewZoneService(zoneRepo)
	rackService := service.NewRackService(rackRepo)
	shelfService := service.NewShelfService(shelfRepo)
	productService := service.NewProductService(productRepo)
	stockService := service.NewStockService(stockRepo)

	authHandler := handler.NewAuthHandler(authService)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	rackHandler := handler.NewRackHandler(rackService)
	shelfHandler := handler.NewShelfHandler(shelfService)
	productHandler := handler.NewProductHandler(productService)
	stockHandler := handler.NewStockHandler(stockService)
	healthHandler := handler.NewHealthHandler(db)

	r := chi.NewRouter()
	r.Use(handler.LoggingMiddleware)
	r.Use(handler.TimeoutMiddleware(5 * time.Second))

	r.Get("/healthz", healthHandler.HealthzHandler)
	r.Post("/register", authHandler.RegisterHandler)
	r.Post("/login", authHandler.LoginHandler)
	r.Post("/refresh", authHandler.RefreshHandler)
	r.Post("/logout", authHandler.LogoutHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(jwtSecret))

		r.Get("/warehouses", warehouseHandler.ListHandler)
		r.Get("/zones", zoneHandler.ListHandler)
		r.Get("/racks", rackHandler.ListHandler)
		r.Get("/shelves", shelfHandler.ListHandler)
		r.Get("/products", productHandler.ListHandler)
		r.Get("/stock", stockHandler.GetByShelfHandler)
		r.Post("/stock/move", stockHandler.MoveHandler)
		r.Post("/stock/add", stockHandler.AddHandler)

		r.Group(func(r chi.Router) {
			r.Use(handler.RequireRole("manager", userRepo))

			r.Post("/warehouses", warehouseHandler.CreateHandler)
			r.Post("/zones", zoneHandler.CreateHandler)
			r.Post("/racks", rackHandler.CreateHandler)
			r.Post("/shelves", shelfHandler.CreateHandler)
			r.Post("/products", productHandler.CreateHandler)
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		slog.Info("server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start: ", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown: ", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
