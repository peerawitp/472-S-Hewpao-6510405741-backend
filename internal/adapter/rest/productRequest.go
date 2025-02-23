package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type ProductRequestHandler interface {
	CreateProductRequest(c *fiber.Ctx) error
	GetDetailByID(c *fiber.Ctx) error
}

type productRequestHandler struct {
	service usecase.ProductRequestUsecase
}

func NewProductRequestHandler(service usecase.ProductRequestUsecase) ProductRequestHandler {
	return &productRequestHandler{
		service: service,
	}
}

func (pr *productRequestHandler) CreateProductRequest(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)

	userId, _ := claims["id"].(string)

	fileHeaders, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fileReaders, files, err := util.FileManage(fileHeaders, "images", 15)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	req := dto.CreateProductRequestRequestDTO{}
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validationErr := util.ValidateStruct(req)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErr.Error,
		})
	}

	productRequest := domain.ProductRequest{
		Name:     req.Name,
		Desc:     req.Desc,
		Budget:   req.Budget,
		Quantity: req.Quantity,
		Category: req.Category,
		Offers:   []domain.Offer{},
		Images:   []string{},
		UserID:   &userId,
	}

	err = pr.service.CreateProductRequest(&productRequest, files, fileReaders)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var deletedAt *time.Time
	if productRequest.DeletedAt.Valid {
		deletedAt = &productRequest.DeletedAt.Time
	}

	res := dto.CreateProductRequestResponseDTO{
		Name:     productRequest.Name,
		Desc:     productRequest.Desc,
		Images:   productRequest.Images,
		Budget:   productRequest.Budget,
		Quantity: productRequest.Quantity,
		Category: productRequest.Category,

		UserID: productRequest.UserID,
		User:   productRequest.User,
		Offers: productRequest.Offers,

		CreatedAt: productRequest.CreatedAt,
		UpdatedAt: productRequest.UpdatedAt,
		DeletedAt: deletedAt,
	}

	return c.JSON(fiber.Map{
		"message":         "Product request created sucessfully",
		"product-request": res,
	})
}

func (pr *productRequestHandler) GetDetailByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	productRequest, err := pr.service.GetDetailByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":         "success",
		"product-request": productRequest,
	})
}
