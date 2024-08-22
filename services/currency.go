package services

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type CurrencyService struct {
	Client  *resty.Client
	APIKey  string
	BaseURL string
}

func NewCurrencyService() *CurrencyService {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("EXCHANGERATE_API_KEY")
	baseURL := "https://v6.exchangerate-api.com/v6"

	client := resty.New()

	return &CurrencyService{
		Client:  client,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

func (s *CurrencyService) ConvertCurrency(from string, to string, amount float64) (float64, error) {
	response := struct {
		Result         string  `json:"result"`
		ConversionRate float64 `json:"conversion_rate"`
		Error          struct {
			Code    string `json:"error-type"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}{}

	endpoint := fmt.Sprintf("%s/%s/pair/%s/%s", s.BaseURL, s.APIKey, from, to)
	log.Printf("Requesting Exchangerate-API.com: %s", endpoint)

	resp, err := s.Client.R().SetResult(&response).Get(endpoint)

	if err != nil {
		log.Printf("Error making API request: %v", err)
		return 0, err
	}

	log.Printf("API response: %s", resp.String())

	if response.Result != "success" {
		log.Printf("Currency conversion failed with error: %v", response.Error)
		return 0, fmt.Errorf("Currency conversion failed: %s", response.Error.Message)
	}

	// Gunakan langsung ConversionRate dari respons
	log.Printf("Conversion rate: %f, Amount before conversion: %f", response.ConversionRate, amount)
	convertedAmount := response.ConversionRate * amount
	log.Printf("Converted amount: %f", convertedAmount)

	return convertedAmount, nil
}
