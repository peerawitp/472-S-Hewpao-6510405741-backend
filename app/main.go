package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hewpao/hewpao-backend/bootstrap"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/ctx"
	"github.com/hewpao/hewpao-backend/internal/adapter/gorm"
	"github.com/hewpao/hewpao-backend/internal/adapter/middleware"
	"github.com/hewpao/hewpao-backend/internal/adapter/oauth"
	"github.com/hewpao/hewpao-backend/internal/adapter/rest"
	"github.com/hewpao/hewpao-backend/internal/adapter/s3"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/usecase"
)

func main() {
	app := fiber.New()
	cfg := config.NewConfig()
	db := bootstrap.NewDB(&cfg)
	ctx := ctx.ProvideContext()
	minio := bootstrap.ProvideMinIOClient(ctx, &cfg)

	app.Use(logger.New())

	oauthRepoFactory := repository.NewOAuthRepositoryFactory()
	oauthRepoFactory.Register("google", oauth.NewGoogleOAuthRepository(&cfg))

	minioRepo := s3.ProvideMinIOS3Repository(minio, &cfg)

	userRepo := gorm.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	authUsecase := usecase.NewAuthUsecase(userRepo, &oauthRepoFactory, &cfg, minioRepo, ctx)
	authHandler := rest.NewAuthHandler(authUsecase)

	productRequestRepo := gorm.NewProductRequestGormRepo(db)
	productRequestUsecase := usecase.NewProductRequestService(productRequestRepo, minioRepo, ctx)
	productRequestHandler := rest.NewProductRequestHandler(productRequestUsecase)

	verifcationUsecase := usecase.NewVerificationService(minioRepo, ctx, cfg, userRepo)
	verifcationHandler := rest.NewVerificationHandler(verifcationUsecase)

	offerRepo := gorm.NewOfferGormRepo(db)
	offerUsecase := usecase.NewOfferService(offerRepo, productRequestRepo, userRepo, ctx)
	offerHandler := rest.NewOfferHandler(offerUsecase)

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

	productRequestRoute := app.Group("/product-requests", middleware.AuthMiddleware(&cfg))
	productRequestRoute.Post("/", productRequestHandler.CreateProductRequest)
	productRequestRoute.Get("/get", productRequestHandler.GetPaginatedProductRequests)
	productRequestRoute.Get("/get/:id", productRequestHandler.GetDetailByID)
	productRequestRoute.Get("/get-buyer", productRequestHandler.GetBuyerProductRequestsByUserID)

	verifyRoute := app.Group("/verify", middleware.AuthMiddleware(&cfg))
	verifyRoute.Post("/", verifcationHandler.VerifyWithKYC)
	verifyRoute.Get("/:email", verifcationHandler.GetVerificationInfo)
	verifyRoute.Post("/set/:email", verifcationHandler.UpdateVerificationInfo)

	offerRoute := app.Group("/offers", middleware.AuthMiddleware(&cfg))
	offerRoute.Post("/", offerHandler.CreateOffer)

	app.Listen(":9090")
}
