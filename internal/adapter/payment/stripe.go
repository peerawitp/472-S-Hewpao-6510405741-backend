package payment

import (
	"context"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type StripePaymentRepository struct {
	cfg *config.Config
}

func NewStripePaymentRepository(cfg *config.Config) repository.PaymentRepository {
	return &StripePaymentRepository{
		cfg: cfg,
	}
}

func (r *StripePaymentRepository) CreatePayment(ctx context.Context, pr *domain.ProductRequest) (*dto.CreatePaymentResponseDTO, error) {
	stripe.Key = r.cfg.StripeSecretKey
	uri := "https://hewpao-fe.peerawitp.me/my-product/"

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(uri),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyTHB)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(pr.Name),
					},
					UnitAmount: stripe.Int64(int64(pr.Budget * 100)),
				},
				Quantity: stripe.Int64(int64(pr.Quantity)),
			},
			{
				// fee
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyTHB)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Service Fee"),
					},
					UnitAmount: stripe.Int64(int64(pr.Budget * 0.05 * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		PaymentMethodTypes: []*string{
			stripe.String("card"),
			stripe.String("promptpay"),
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		return nil, err
	}

	res := &dto.CreatePaymentResponseDTO{
		PaymentID:  result.ID,
		PaymentURL: result.URL,
		CreatedAt:  result.Created,
		ExpiredAt:  result.ExpiresAt,
	}

	return res, nil
}
