package webhook

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/types"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeWebhookHandler interface {
	WebhookPost(c *fiber.Ctx) error
}

type stripeWebhookHandler struct {
	cfg       *config.Config
	service   usecase.CheckoutUsecase
	txService usecase.TransactionUseCase
	prService usecase.ProductRequestUsecase
}

func NewStripeWebhookHandler(cfg *config.Config, service usecase.CheckoutUsecase, txService usecase.TransactionUseCase, prService usecase.ProductRequestUsecase) StripeWebhookHandler {
	return &stripeWebhookHandler{
		cfg:       cfg,
		service:   service,
		txService: txService,
		prService: prService,
	}
}

func (s *stripeWebhookHandler) WebhookPost(c *fiber.Ctx) error {
	stripe.Key = s.cfg.StripeSecretKey

	payload := c.Body()
	stripeSignature := c.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, stripeSignature, s.cfg.StripeWebhookSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	switch event.Type {
	case "checkout.session.completed":
		var paymentData stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &paymentData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = s.service.UpdateTransactionStatus(c.Context(), paymentData.ID, types.PaymentSuccess)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		transaction, err := s.txService.GetTransactionByThirdPartyPaymentID(c.Context(), paymentData.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		prID := int(*transaction.ProductRequestID)
		updateErr := s.prService.UpdateProductRequestStatusAfterPaid(prID)
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": updateErr.Error(),
			})
		}

	case "checkout.session.expired":
		var paymentData stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &paymentData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = s.service.UpdateTransactionStatus(c.Context(), paymentData.ID, types.PaymentExpired)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return c.SendString("stripe webhook received 💸!")
}
