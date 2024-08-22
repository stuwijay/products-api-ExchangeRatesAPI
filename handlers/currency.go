package handlers

import (
	"net/http"
	"strconv"
	"warungjwt_postgre2/services"
	"warungjwt_postgre2/utils"

	"github.com/labstack/echo/v4"
)

var currencyService = services.NewCurrencyService()

func ConvertCurrencyHandler(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	amountStr := c.QueryParam("amount")

	if from == "" || to == "" || amountStr == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Please provide 'from', 'to', and 'amount' query parameters"))
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid amount value"))
	}

	convertedAmount, err := currencyService.ConvertCurrency(from, to, amount)
	if err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to convert currency"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"from":             from,
		"to":               to,
		"original_amount":  amount,
		"converted_amount": convertedAmount,
	})
}
