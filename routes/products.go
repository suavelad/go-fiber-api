package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/suavelad/go-fibre-api/database"
	"github.com/suavelad/go-fibre-api/models"
)

type Product struct {
	Id           uint   `json:"id"`
	Name         string `validate:"required" json:"name"`
	SerialNumber string `validate:"required" json:"serial_number"`
}

func CreateResponseProduct(product models.Product) Product {
	return Product{
		Id:           product.Id,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
	}
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	// Validate the struct using the validator
	var requestProduct Product
	if err := c.BodyParser(&requestProduct); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := validate.Struct(requestProduct); err != nil {
		// Return validation error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_error": err.Error()})
	}

	product = models.Product{
		Name:         requestProduct.Name,
		SerialNumber: requestProduct.SerialNumber,
	}
	database.DB.Db.Create(&product)
	responseProduct := CreateResponseProduct(product)

	return c.Status(201).JSON(responseProduct)
}

func GetProducts(c *fiber.Ctx) error {
	products := []models.Product{}

	database.DB.Db.Find(&products)
	productsResponse := []Product{}

	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		productsResponse = append(productsResponse, responseProduct)
	}

	return c.Status(200).JSON(productsResponse)
}

func findProduct(id int, product *models.Product) error {

	database.DB.Db.Find((product), "id=?", id)
	if product.Id == 0 {
		return errors.New("product does not exist")
	}
	return nil
}

func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err)
	}

	responseProduct := CreateResponseProduct(product)
	return c.Status(200).JSON(responseProduct)

}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err)
	}

	type UpdateProduct struct {
		Name         string `validate:"required" json:"name"`
		SerialNumber string `validate: "required" json:"serial_number"`
	}

	var updateData UpdateProduct

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := validate.Struct(updateData); err != nil {
		
		// Return validation error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_error": err.Error()})
	}
	product.Name = updateData.Name
	product.SerialNumber = updateData.SerialNumber

	database.DB.Db.Save(&product)

	productResponse := CreateResponseProduct(product)

	return c.Status(200).JSON(productResponse)

}

func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err)
	}
	if err := database.DB.Db.Delete(&product).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully deleted")
}
