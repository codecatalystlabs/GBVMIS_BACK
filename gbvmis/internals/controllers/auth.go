package controllers

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Login godoc
//	@Summary		Login a police officer
//	@Description	Authenticate a police officer using email or username and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		object{identifier=string,password=string}	true	"Login credentials"
//	@Success		200			{object}	map[string]string							"Returns access and refresh tokens"
//	@Failure		400			{object}	map[string]string							"Invalid input"
//	@Failure		401			{object}	map[string]string							"Invalid credentials"
//	@Failure		500			{object}	map[string]string							"Token generation error"
//	@Router			/login [post]
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

// RefreshToken godoc
//	@Summary		Refresh JWT tokens
//	@Description	Generates new access and refresh tokens using a valid refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			refresh_token	body		object{refresh_token=string}	true	"Refresh token"
//	@Success		200				{object}	map[string]string				"Returns new access and refresh tokens"
//	@Failure		400				{object}	map[string]string				"Invalid input"
//	@Failure		401				{object}	map[string]string				"Invalid or expired refresh token"
//	@Failure		500				{object}	map[string]string				"Token generation error"
//	@Router			/refresh-token [post]
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
