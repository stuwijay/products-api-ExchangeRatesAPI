package main

import (
	"net/http"
	"warungjwt_postgre2/config"
	"warungjwt_postgre2/handlers"
	"warungjwt_postgre2/middlewares"
	"warungjwt_postgre2/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	// Default status code
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	// Type assertion to check if the error is an APIError
	if apiErr, ok := err.(*utils.APIError); ok {
		code = apiErr.Code
		message = apiErr.Message
	} else if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Send the error response
	c.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
	})
}

func main() {
	config.InitDatabase()

	e := echo.New()

	// Custom HTTPErrorHandler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Middleware logger
	e.Use(middleware.Logger())

	// Middleware recovery (optional, handles panics)
	e.Use(middleware.Recover())

	// Public routes
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	// Currency conversion route
	e.GET("/convert", handlers.ConvertCurrencyHandler)

	// Protected routes - Role-based access
	products := e.Group("/products")
	products.Use(middlewares.IsAuthorized("admin")) // Middleware to allow only admin for POST, PUT, DELETE
	products.GET("", handlers.GetProducts)          // Accessible by both admin and staff
	products.GET("/:id", handlers.GetProduct)       // Accessible by both admin and staff
	products.POST("", handlers.CreateProduct)
	products.PUT("/:id", handlers.UpdateProduct)
	products.DELETE("/:id", handlers.DeleteProduct)

	// Route for staff (only GET access)
	staffProducts := e.Group("/staff/products")
	staffProducts.Use(middlewares.IsAuthorized("staff")) // Middleware to allow only staff with GET
	staffProducts.GET("", handlers.GetProducts)
	staffProducts.GET("/:id", handlers.GetProduct)

	e.Logger.Fatal(e.Start(":1323"))
}
