package handlers

import (
	"log"
	"net/http"
	"strconv"
	"warungjwt_postgre2/config"
	"warungjwt_postgre2/models"
	"warungjwt_postgre2/utils"

	"github.com/labstack/echo/v4"
)

func GetProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid product ID"))
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Produk tidak ditemukan"))
	}

	// Log untuk memastikan harga dari database benar
	log.Printf("Harga produk dari database: %f", product.Price)

	// Menambahkan konversi mata uang
	toCurrency := c.QueryParam("currency")
	if toCurrency != "" {
		convertedPrice, err := currencyService.ConvertCurrency("USD", toCurrency, product.Price)
		if err != nil {
			return utils.HandleError(c, utils.NewInternalError("Failed to convert currency"))
		}
		log.Printf("Converted product price: %f", convertedPrice)
		product.Price = convertedPrice
	}

	return c.JSON(http.StatusOK, product)
}

func GetProducts(c echo.Context) error {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError(err.Error()))
	}

	// Opsi untuk mengkonversi harga semua produk ke mata uang yang diinginkan
	toCurrency := c.QueryParam("currency")
	if toCurrency != "" {
		for i := range products {
			convertedPrice, err := currencyService.ConvertCurrency("USD", toCurrency, products[i].Price)
			if err != nil {
				return utils.HandleError(c, utils.NewInternalError("Failed to convert currency for product ID: "+strconv.Itoa(int(products[i].ID))))
			}
			products[i].Price = convertedPrice
		}
	}

	return c.JSON(http.StatusOK, products)
}

func CreateProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid request payload"))
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to create product"))
	}

	return c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid product ID"))
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Produk tidak ditemukan"))
	}

	if err := c.Bind(&product); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid request payload"))
	}

	if err := config.DB.Save(&product).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to update product"))
	}

	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid product ID"))
	}

	// Pastikan bahwa penghapusan dilakukan dengan primary key yang benar
	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to delete product"))
	}

	// Cek apakah produk benar-benar terhapus
	rowsAffected := config.DB.Delete(&models.Product{}, id).RowsAffected
	if rowsAffected == 0 {
		return utils.HandleError(c, utils.NewNotFoundError("Product not found or already deleted"))
	}

	return c.NoContent(http.StatusNoContent)
}
