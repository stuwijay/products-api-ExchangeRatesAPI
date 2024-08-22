package handlers

import (
	"net/http"
	"time"
	"warungjwt_postgre2/config"
	"warungjwt_postgre2/models"
	"warungjwt_postgre2/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Generate JWT token
func generateJWT(email, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expires in 72 hours

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}

func Register(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	fullName := c.FormValue("full_name")
	role := c.FormValue("role")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
		Role:     role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Gagal menyimpan pengguna"))
	}

	// Generate JWT
	token, err := generateJWT(user.Email, user.Role)
	if err != nil {
		return utils.HandleError(c, utils.NewInternalError("Gagal menghasilkan token"))
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Pengguna berhasil didaftarkan",
		"token":   token,
	})
}

func Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Pengguna tidak ditemukan"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Kata sandi salah"))
	}

	// Generate JWT
	token, err := generateJWT(user.Email, user.Role)
	if err != nil {
		return utils.HandleError(c, utils.NewInternalError("Gagal menghasilkan token"))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login berhasil",
		"token":   token,
	})
}
