package routes

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/suavelad/go-fibre-api/database"
	"github.com/suavelad/go-fibre-api/models"
)

type Order struct {
	Id        uint      `json:"id"  gorm:"primaryKey"`
	Title     string    `validate:"required" json:"title"`
	Product   Product   `validate:"required" json:"product"`
	User      User      `validate:"required" json:"user"`
	CreatedAt time.Time `json:"order_date "`
}

var validate = validator.New()

func CreateOrderResponse(order models.Order, user User, product Product) Order {
	return Order{
		Id:        order.Id,
		Title:     order.Title,
		Product:   product,
		User:      user,
		CreatedAt: order.CreatedAt,
	}
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	var request Order

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// Validate the struct using the validator
	if err := validate.Struct(request); err != nil {
		// Return validation error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_error": err.Error()})
	}
	order = models.Order{
		Title:        request.Title,
		ProductRefer: int(request.Product.Id),
		UserRefer:    int(request.User.Id),
	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(400).JSON(err.Error)
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		// return c.Status(400).JSON(err.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Db.Create(&order)
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)

	orderResponse := CreateOrderResponse(order, responseUser, responseProduct)
	return c.Status(200).JSON(orderResponse)

}

func findOrder(id uint, order *models.Order) error {
	database.DB.Db.Find(&order, "id = ?", id)
	if order.Id == 0 {
		return errors.New("order does not exist")
	}
	return nil
}

func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.New("invalid id")
	}
	var order models.Order
	if err := findOrder(uint(id), &order); err != nil {
		return c.Status(400).JSON(err.Error)
	}
	createOrderResponse := CreateOrderResponse(order, CreateResponseUser(order.User), CreateResponseProduct(order.Product))

	return c.Status(200).JSON(createOrderResponse)

}

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Db.Find(&orders)
	ordersResponse := []Order{}

	for _, o := range orders {
		ordersResponse = append(ordersResponse, CreateOrderResponse(o, CreateResponseUser(o.User), CreateResponseProduct(o.Product)))
	}
	return c.Status(200).JSON(ordersResponse)

}

func DeleteOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.New("invalid id")
	}
	var order models.Order
	if err := findOrder(uint(id), &order); err != nil {
		return c.Status(400).JSON(err.Error)
	}

	database.DB.Db.Delete(&order)
	return c.Status(200).JSON("Deleted Successfully")

}

func UpdateOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.New("invalid id")
	}
	var order models.Order

	if err := findOrder(uint(id), &order); err != nil {
		return c.Status(400).JSON(err.Error)
	}

	type UpdateOrder struct {
		Title   string `validate:"required" json:"title"`
		Product uint   `validate:"required" json:"product"`
		User    uint   `validate:"required" json:"user"`
	}
	var updateData UpdateOrder

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	
	// Validate the struct using the validator
	if err := validate.Struct(updateData); err != nil {
		// Return validation error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_error": err.Error()})
	}
	
	order.Title = updateData.Title
	order.UserRefer = int(updateData.User)
	order.ProductRefer = int(updateData.Product)

	database.DB.Db.Save(&order)
	return c.Status(200).JSON(order)
}
