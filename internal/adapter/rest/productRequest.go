package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/util"
)

type ProductRequestHandler interface {
	CreateProductRequest(c *fiber.Ctx) error
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
	fileHeaders, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	fileReaders, files, err := util.FileManage(fileHeaders, "images", 15)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	req := dto.CreateProductRequestRequestDTO{}
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	validationErr := util.ValidateStruct(req)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).SendString(validationErr.Error)
	}

	productRequest := domain.ProductRequest{
		Name:     req.Name,
		Desc:     req.Desc,
		Budget:   req.Budget,
		Quantity: req.Quantity,
		Category: req.Category,
		Offers:   []domain.Offer{},
		Images:   []string{},
	}

	err = pr.service.CreateProductRequest(&productRequest, files, fileReaders)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
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
