package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Wizardsmile1412/hospital-middleware/internal/client"
	"github.com/Wizardsmile1412/hospital-middleware/internal/config"
	"github.com/Wizardsmile1412/hospital-middleware/internal/handler"
	"github.com/Wizardsmile1412/hospital-middleware/internal/repository"
	"github.com/Wizardsmile1412/hospital-middleware/internal/routes"
	"github.com/Wizardsmile1412/hospital-middleware/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	runMigrations(cfg)

	db := connectDB(cfg)
	defer db.Close()

	// Repositories
	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	// Clients
	hospitalClient := client.NewHospitalAClient(cfg.HospitalAAPIURL)

	// Services
	staffSvc := service.NewStaffService(staffRepo, cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	patientSvc := service.NewPatientService(patientRepo, hospitalClient)

	// Handlers
	staffHandler := handler.NewStaffHandler(staffSvc, cfg)
	patientHandler := handler.NewPatientHandler(patientSvc)

	// Router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	routes.Register(r, staffHandler, patientHandler, cfg)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func connectDB(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected")
	return pool
}

func runMigrations(cfg *config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied")
}
