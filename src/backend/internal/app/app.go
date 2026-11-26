package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/TwiLightDM/lab-track/internal/config"
	_ "github.com/TwiLightDM/lab-track/internal/docs"
	"github.com/TwiLightDM/lab-track/internal/middlewares"
	"github.com/TwiLightDM/lab-track/internal/repository"
	"github.com/TwiLightDM/lab-track/internal/service"
	httpapi "github.com/TwiLightDM/lab-track/internal/transport/v1/http"
	"github.com/TwiLightDM/lab-track/package/databases"
	"github.com/TwiLightDM/lab-track/package/tools"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Run(cfg *config.Config) {
	db, err := databases.InitDB(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}

	jwtService := tools.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessDuration,
		cfg.JWT.RefreshDuration,
	)

	validatorService := tools.NewValidatorService()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, jwtService)
	userHandler := httpapi.NewUserHandler(userService, validatorService)

	mw := middlewares.NewMWManager(jwtService, validatorService)

	registerRoutes(e, mw, userHandler)

	go func() {
		for attempt := 1; attempt <= 5; attempt++ {
			if err = e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
				e.Logger.Errorf("Failed to start server: %v attempt: %d\n", err, attempt)
				time.Sleep(5 * time.Second)
			} else {
				return
			}
		}
		log.Fatal("Failed to start server")
	}()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown server: %v\n", err)
	}

	if err = sqlDB.Close(); err != nil {
		log.Printf("Failed to close database connection: %v\n", err)
		return
	}
	log.Println("Shutdown complete")
}

func registerRoutes(
	e *echo.Echo,
	mw *middlewares.MWManager,
	userHandler *httpapi.UserHandler,
) {
	v1 := e.Group("api/v1")

	public := v1.Group("/public")
	users := v1.Group("/users", mw.BearerAuth())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	public.POST("/login", userHandler.Login)
	public.POST("/refresh", userHandler.RefreshTokens, mw.BearerAuth())

	users.POST("", userHandler.CreateUser, mw.AdminCheck())
	users.GET("", userHandler.GetAllUsers)
	users.GET("/:id", userHandler.GetUserById)
	users.PATCH("/:id", userHandler.PatchUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
