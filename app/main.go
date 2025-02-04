package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hewpao/hewpao-backend/bootstrap"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/internal/adapter/gorm"
	"github.com/hewpao/hewpao-backend/internal/adapter/oauth"
	"github.com/hewpao/hewpao-backend/internal/adapter/rest"
	"github.com/hewpao/hewpao-backend/usecase"
)

func main() {
	app := fiber.New()
	cfg := config.NewConfig()
	db := bootstrap.NewDB(&cfg)

	app.Use(logger.New())

	oauthFactory := oauth.NewOAuthServiceFactory()
	oauthFactory.Register("google", oauth.NewGoogleOAuthService(&cfg))

	userRepo := gorm.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	authUsecase := usecase.NewAuthUsecase(userRepo, &oauthFactory, &cfg)
	authHandler := rest.NewAuthHandler(authUsecase)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hewpao is running 🚀")
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		user, err := userUsecase.GetUserByID(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		return c.JSON(user)
	})

	authRoute := app.Group("/auth")
	authRoute.Post("/login", authHandler.LoginWithCredentials)
	authRoute.Post("/login/oauth", authHandler.LoginWithOAuth)
	authRoute.Post("/register", authHandler.Register)

	app.Listen(":9090")
}
