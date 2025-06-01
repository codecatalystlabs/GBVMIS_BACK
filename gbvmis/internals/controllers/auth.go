package controllers

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var creds struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}

		if err := c.BodyParser(&creds); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		var officer models.PoliceOfficer
		if err := db.Preload("Roles").
			Where("email = ? OR username = ?", creds.Identifier, creds.Identifier).
			First(&officer).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if !utils.CheckPassword(officer.Password, creds.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		var roleNames []string
		for _, r := range officer.Roles {
			roleNames = append(roleNames, r.Name)
		}

		accessToken, refreshToken, err := utils.GenerateTokens(officer.ID, officer.Email, roleNames)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create tokens"})
		}

		return c.JSON(fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := jwt.ParseWithClaims(body.RefreshToken, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.RefreshKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired refresh token"})
	}

	claims, ok := token.Claims.(*utils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Create new access & refresh tokens
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Email, claims.Roles)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tokens"})
	}

	return c.JSON(fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
