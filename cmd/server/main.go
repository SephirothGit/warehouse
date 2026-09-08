package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/SephirothGit/warehouse/internal/handler"
	"github.com/SephirothGit/warehouse/internal/repository"
	"github.com/SephirothGit/warehouse/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("unable to open db connection", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("unable to connect to db: ", err)
	}
	log.Println("connected to database successfully")

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
}