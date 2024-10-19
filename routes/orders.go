package routes

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suavelad/go-fibre-api/database"
	"github.com/suavelad/go-fibre-api/models"
)

type Order struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Product   Product   `json:"product"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"order_date "`
}

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
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(400).JSON(err.Error)
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(400).JSON(err.Error)
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
		Title   string `json:"title"`
		Product uint   `json:"product"`
		User    uint   `json:"user"`
	}
	var updateData UpdateOrder
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	order.Title = updateData.Title
	order.UserRefer = int(updateData.User)
	order.ProductRefer = int(updateData.Product)

	database.DB.Db.Save(&order)
	return c.Status(200).JSON(order)
}
