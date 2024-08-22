
# Warung Postgre Project

## 1. Pengantar Project
Project awal ini bernama `warung_postgre`, yang bertujuan untuk mengelola produk dan pengguna dalam sebuah sistem berbasis API menggunakan Go dan PostgreSQL. 

## 2. Penambahan Kolom `price` pada Tabel `products`
Untuk mendukung fitur konversi harga produk dalam berbagai mata uang, kita menambahkan kolom `price` pada tabel `products`.

### **DDL.sql**
Berikut adalah perubahan yang dilakukan pada tabel `products`:

```sql
ALTER TABLE products
ADD COLUMN price NUMERIC(10, 2) NOT NULL DEFAULT 0.00;
```

### **DML.sql**
Setelah menambahkan kolom `price`, berikut adalah data yang diperbarui dan data baru yang ditambahkan ke dalam tabel:

```sql
-- Memperbarui harga untuk produk yang sudah ada
UPDATE products SET price = 100.00 WHERE code = 'P001';
UPDATE products SET price = 150.00 WHERE code = 'P002';
UPDATE products SET price = 200.00 WHERE code = 'P003';
UPDATE products SET price = 250.00 WHERE code = 'P099';
UPDATE products SET price = 300.00 WHERE code = 'P098';
UPDATE products SET price = 350.00 WHERE code = 'P077';

-- Menambahkan produk baru dengan kolom 'price' yang diisi
INSERT INTO products (name, code, stock, description, status, price) VALUES
('Product 4', 'P004', 20, 'Description for Product 4', 'active', 120.00),
('Product 5', 'P005', 25, 'Description for Product 5', 'active', 180.00);
```

`notes : sesuaikan dengan data product yang ada pada masing-masing database Anda ya.` 

## 3. Implementasi Exchangerate-API

Untuk menambahkan fitur konversi mata uang secara real-time, kita menggunakan Exchangerate-API. Berikut adalah langkah-langkah implementasinya:

### **Step 1: Mendapatkan API Key**
- Daftar di [Exchangerate-API.com](https://www.exchangerate-api.com/) dan dapatkan API key Anda.
- Tambahkan API key tersebut ke dalam file `.env`:
    ```env
    EXCHANGERATE_API_KEY=your_exchangerate_api_key
    ```

### **Step 2: Membuat Service untuk Mengakses API**
Buat file `services/currency.go` untuk menangani koneksi ke Exchangerate-API:

```go
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
        Result          string  `json:"result"`
        ConversionRate  float64 `json:"conversion_rate"`
        Error           struct {
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
```

### **Step 3: Mengintegrasikan dengan API Produk**
Implementasikan konversi mata uang ke dalam handler produk (`handlers/product.go`):

```go
package handlers

import (
    "net/http"
    "strconv"
    "warungjwt_postgre2/config"
    "warungjwt_postgre2/models"
    "warungjwt_postgre2/services"
    "warungjwt_postgre2/utils"

    "github.com/labstack/echo/v4"
)

var currencyService = services.NewCurrencyService()

func GetProduct(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return utils.HandleError(c, utils.NewBadRequestError("Invalid product ID"))
    }

    var product models.Product
    if err := config.DB.First(&product, id).Error; err != nil {
        return utils.HandleError(c, utils.NewNotFoundError("Produk tidak ditemukan"))
    }

    toCurrency := c.QueryParam("currency")
    if toCurrency != "" {
        convertedPrice, err := currencyService.ConvertCurrency("USD", toCurrency, product.Price)
        if err != nil {
            return utils.HandleError(c, utils.NewInternalError("Failed to convert currency"))
        }
        product.Price = convertedPrice
    }

    return c.JSON(http.StatusOK, product)
}

func GetProducts(c echo.Context) error {
    var products []models.Product
    if err := config.DB.Find(&products).Error; err != nil {
        return utils.HandleError(c, utils.NewInternalError(err.Error()))
    }

    toCurrency := c.QueryParam("currency")
    if toCurrency != "" {
        for i := range products {
            convertedPrice, err := currencyService.ConvertCurrency("USD", toCurrency, products[i].Price)
            if err != nil {
                return utils.HandleError(c, utils.NewInternalError("Failed to convert currency for product ID: " + strconv.Itoa(int(products[i].ID))))
            }
            products[i].Price = convertedPrice
        }
    }

    return c.JSON(http.StatusOK, products)
}
```

## 4. Dokumentasi Jalankan Endpoint

Setelah semua konfigurasi selesai, Anda bisa mengakses API untuk mendapatkan produk dengan harga yang dikonversi ke mata uang lain:

### **Endpoint untuk Mendapatkan Produk Berdasarkan ID**
- **URL:** `GET http://localhost:1323/products/:id`
- **Parameter Query:** 
  - `currency`: kode mata uang yang diinginkan, misalnya `EUR`, `IDR`, dll.
- **Contoh Request:**
  ```
  GET http://localhost:1323/products/1?currency=EUR
  ```
- **Contoh Response:**
  ```json
  {
    "id": 1,
    "name": "Product 1",
    "code": "P001",
    "stock": 10,
    "description": "Description for Product 1",
    "status": "active",
    "price": 85.50  // Harga dalam EUR setelah konversi
  }
  ```

Dengan mengikuti langkah-langkah ini, Anda dapat menambahkan dan mengonversi harga produk secara real-time menggunakan Exchangerate-API.
