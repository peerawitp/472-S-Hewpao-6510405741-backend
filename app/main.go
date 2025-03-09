package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hewpao/hewpao-backend/bootstrap"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/ctx"
	"github.com/hewpao/hewpao-backend/internal/adapter/ekyc"
	"github.com/hewpao/hewpao-backend/internal/adapter/gorm"
	"github.com/hewpao/hewpao-backend/internal/adapter/middleware"
	"github.com/hewpao/hewpao-backend/internal/adapter/notitype"
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
	httpCli := &http.Client{}

	app.Use(logger.New())

	offerRepo := gorm.NewOfferGormRepo(db)

	oauthRepoFactory := repository.NewOAuthRepositoryFactory()
	oauthRepoFactory.Register("google", oauth.NewGoogleOAuthRepository(&cfg))

	paymentRepoFactory := repository.NewPaymentRepositoryFactory()
	paymentRepoFactory.Register("stripe", payment.NewStripePaymentRepository(&cfg))

	minioRepo := s3.ProvideMinIOS3Repository(minio, &cfg)

	userRepo := gorm.NewUserGormRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	notificationRepoFactory := repository.NewNotificationRepositoryFactory()
	emailRepo, err := notitype.NewEmailNotificationRepo(message, &cfg)
	if err != nil {
		log.Panic(err)
	}
	logRepo, err := notitype.NewTestNotificationRepo(&cfg)
	if err != nil {
		log.Panic(err)
	}

	notificationRepoFactory.Register("email", emailRepo)
	notificationRepoFactory.Register("log", logRepo)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepoFactory, userRepo, ctx, message, &cfg, offerRepo)

	authUsecase := usecase.NewAuthUsecase(userRepo, &oauthRepoFactory, &cfg, minioRepo, ctx)
	authHandler := rest.NewAuthHandler(authUsecase)

	chatRepo := gorm.NewChatRepo(db)
	chatUsecase := usecase.NewChatService(chatRepo)
	chatHandler := rest.NewChatHandler(chatUsecase)

	productRequestRepo := gorm.NewProductRequestGormRepo(db)
	productRequestUsecase := usecase.NewProductRequestService(productRequestRepo, minioRepo, ctx, offerRepo, userRepo, chatRepo, &cfg, message)
	productRequestHandler := rest.NewProductRequestHandler(productRequestUsecase, notificationUsecase)

	transactionRepo := gorm.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionService(transactionRepo)
	transactionHandler := rest.NewTransactionHandler(*transactionUsecase)

	verificationRepo := gorm.NewVerificationGormRepo(db)

	ekycRepoFactory := repository.NewEKYCRepositoryFactory()
	ekycRepoFactory.Register("iapp", ekyc.NewIappVerificationRepo(&cfg, httpCli))

	verifcationUsecase := usecase.NewVerificationService(minioRepo, ctx, cfg, userRepo, verificationRepo, ekycRepoFactory)
	verifcationHandler := rest.NewVerificationHandler(verifcationUsecase)

	offerUsecase := usecase.NewOfferService(offerRepo, productRequestRepo, userRepo, ctx)
	offerHandler := rest.NewOfferHandler(offerUsecase)

	checkoutUsecase := usecase.NewCheckoutUsecase(userRepo, productRequestRepo, transactionRepo, &paymentRepoFactory, &cfg, minioRepo, ctx)
	checkoutHandler := rest.NewCheckoutHandler(checkoutUsecase)

	stripeWebhookHandler := webhook.NewStripeWebhookHandler(&cfg, checkoutUsecase)

	messageRepo := gorm.NewMessageGormRepo(db)
	messageUsecase := usecase.NewMessageService(messageRepo)
	messageHandler := rest.NewMessageHandler(*messageUsecase)

	bankRepo := gorm.NewBankGormRepo(db)

	travlerPayoutAccountRepo := gorm.NewTravelerPayoutAccountGormRepository(db)
	travelerPayoutAccountUsecase := usecase.NewTravelerPayoutAccountService(ctx, travlerPayoutAccountRepo, bankRepo)
	travelerPayoutAccountHandler := rest.NewTravelerPayoutAccountHandler(travelerPayoutAccountUsecase)

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
	verifyRoute.Get("/:verification_id", verifcationHandler.GetVerificationInfo)
	verifyRoute.Post("/set/:email", verifcationHandler.UpdateVerificationInfo)

	offerRoute := app.Group("/offers", middleware.AuthMiddleware(&cfg))
	offerRoute.Post("/", offerHandler.CreateOffer)

	transactionRoute := app.Group("/transactions", middleware.AuthMiddleware(&cfg))
	transactionRoute.Post("/", transactionHandler.CreateTransaction)
	transactionRoute.Get("/:id", transactionHandler.GetTransactionByID)

	checkoutRoute := app.Group("/checkout", middleware.AuthMiddleware(&cfg))
	checkoutRoute.Post("/gateway", checkoutHandler.CheckoutWithPaymentGateway)

	travelerPayoutAccountRoute := app.Group("/payout-account", middleware.AuthMiddleware(&cfg))
	travelerPayoutAccountRoute.Post("/create", travelerPayoutAccountHandler.CreateTravelerPayoutAccount)
	travelerPayoutAccountRoute.Get("/get", travelerPayoutAccountHandler.GetAccountsByUserID)
	travelerPayoutAccountRoute.Get("/get-available-banks", travelerPayoutAccountHandler.GetAllAvailableBank)

	// Webhook route
	webhookRoute := app.Group("/webhook")
	stripeWebhookRoute := webhookRoute.Group("/stripe")
	stripeWebhookRoute.Post("/", stripeWebhookHandler.WebhookPost)

	chatRoute := app.Group("/chat")
	chatRoute.Post("/create", chatHandler.CreateChat)

	messageRoute := app.Group("/message")
	messageRoute.Post("/create", messageHandler.CreateMessage)

	app.Listen(":9090")
}
