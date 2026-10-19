package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SephirothGit/warehouse/internal/auth"
	"github.com/SephirothGit/warehouse/internal/config"
	"github.com/SephirothGit/warehouse/internal/handler"
	"github.com/SephirothGit/warehouse/internal/repository"
	"github.com/SephirothGit/warehouse/internal/service"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DB.URL)
	if err != nil {
		slog.Error("unable to open db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		slog.Error("unable to connect to db", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database successfully")

	jwtSecret := []byte(cfg.Auth.JWTSecret)

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

	limiterCtx, cancelLimiter := context.WithCancel(context.Background())
	defer cancelLimiter()

	authLimiter := handler.NewIPRateLimiter(limiterCtx, rate.Every(cfg.RateLimit.AuthInterval), cfg.RateLimit.AuthBurst)

	r := chi.NewRouter()
	r.Use(handler.SecurityHeadersMiddleware)
	r.Use(handler.LoggingMiddleware)
	r.Use(handler.TimeoutMiddleware(cfg.HTTP.RequestTimeout))

	r.Get("/healthz", healthHandler.HealthzHandler)
	r.With(authLimiter.Middleware).Post("/register", authHandler.RegisterHandler)
	r.With(authLimiter.Middleware).Post("/login", authHandler.LoginHandler)
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
		Addr:              cfg.HTTP.Addr,
		Handler:           r,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	go func() {
		slog.Info("server starting", "addr", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start: ", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
