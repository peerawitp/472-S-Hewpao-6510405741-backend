package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hewpao/hewpao-backend/bootstrap"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/ctx"
	"github.com/hewpao/hewpao-backend/internal/adapter/email"
	"github.com/hewpao/hewpao-backend/internal/adapter/gorm"
	"github.com/hewpao/hewpao-backend/internal/adapter/middleware"
	"github.com/hewpao/hewpao-backend/internal/adapter/oauth"
	"github.com/hewpao/hewpao-backend/internal/adapter/payment"
	"github.com/hewpao/hewpao-backend/internal/adapter/rest"
	"github.com/hewpao/hewpao-backend/internal/adapter/rest/webhook"
	"github.com/hewpao/hewpao-backend/internal/adapter/s3"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/hewpao/hewpao-backend/usecase"
	"gopkg.in/gomail.v2"
)

func main() {
	app := fiber.New()
	cfg := config.NewConfig()
	db := bootstrap.NewDB(&cfg)
	ctx := ctx.ProvideContext()
	minio := bootstrap.ProvideMinIOClient(ctx, &cfg)

	message := gomail.NewMessage()

	app.Use(logger.New())

	oauthRepoFactory := repository.NewOAuthRepositoryFactory()
	oauthRepoFactory.Register("google", oauth.NewGoogleOAuthRepository(&cfg))

	paymentRepoFactory := repository.NewPaymentRepositoryFactory()
	paymentRepoFactory.Register("stripe", payment.NewStripePaymentRepository(&cfg))

	minioRepo := s3.ProvideMinIOS3Repository(minio, &cfg)

	userRepo := gorm.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	gmailRepo, err := email.NewGmailEmailNotificationRepo(message, &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	gmailNotificationUsecase := usecase.NewNotificationUsecase(gmailRepo, userRepo, ctx, message)

	authUsecase := usecase.NewAuthUsecase(userRepo, &oauthRepoFactory, &cfg, minioRepo, ctx)
	authHandler := rest.NewAuthHandler(authUsecase)

	offerRepo := gorm.NewOfferGormRepo(db)
	productRequestRepo := gorm.NewProductRequestGormRepo(db)
	productRequestUsecase := usecase.NewProductRequestService(productRequestRepo, minioRepo, ctx, offerRepo, userRepo, gmailNotificationUsecase, &cfg)
	productRequestHandler := rest.NewProductRequestHandler(productRequestUsecase)

	transactionRepo := gorm.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionService(transactionRepo)
	transactionHandler := rest.NewTransactionHandler(*transactionUsecase)

	verifcationUsecase := usecase.NewVerificationService(minioRepo, ctx, cfg, userRepo)
	verifcationHandler := rest.NewVerificationHandler(verifcationUsecase)

	offerUsecase := usecase.NewOfferService(offerRepo, productRequestRepo, userRepo, ctx)
	offerHandler := rest.NewOfferHandler(offerUsecase)

	checkoutUsecase := usecase.NewCheckoutUsecase(userRepo, productRequestRepo, transactionRepo, &paymentRepoFactory, &cfg, minioRepo, ctx)
	checkoutHandler := rest.NewCheckoutHandler(checkoutUsecase)

	stripeWebhookHandler := webhook.NewStripeWebhookHandler(&cfg, checkoutUsecase)

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
	productRequestRoute.Put("/:id", productRequestHandler.UpdateProductRequest)
	productRequestRoute.Put("/status/:id", productRequestHandler.UpdateProductRequestStatus)
	productRequestRoute.Get("/get", productRequestHandler.GetPaginatedProductRequests)
	productRequestRoute.Get("/get/:id", productRequestHandler.GetDetailByID)
	productRequestRoute.Get("/get-buyer", productRequestHandler.GetBuyerProductRequestsByUserID)

	verifyRoute := app.Group("/verify", middleware.AuthMiddleware(&cfg))
	verifyRoute.Post("/", verifcationHandler.VerifyWithKYC)
	verifyRoute.Get("/:email", verifcationHandler.GetVerificationInfo)
	verifyRoute.Post("/set/:email", verifcationHandler.UpdateVerificationInfo)

	offerRoute := app.Group("/offers", middleware.AuthMiddleware(&cfg))
	offerRoute.Post("/", offerHandler.CreateOffer)

	transactionRoute := app.Group("/transactions", middleware.AuthMiddleware(&cfg))
	transactionRoute.Post("/", transactionHandler.CreateTransaction)
	transactionRoute.Get("/:id", transactionHandler.GetTransactionByID)

	checkoutRoute := app.Group("/checkout", middleware.AuthMiddleware(&cfg))
	checkoutRoute.Post("/gateway", checkoutHandler.CheckoutWithPaymentGateway)

	// Webhook route
	webhookRoute := app.Group("/webhook")
	stripeWebhookRoute := webhookRoute.Group("/stripe")
	stripeWebhookRoute.Post("/", stripeWebhookHandler.WebhookPost)

	app.Listen(":9090")
}
