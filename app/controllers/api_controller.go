package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iwinardhyas/sprint_backend/app/models"
	"github.com/iwinardhyas/sprint_backend/pkg/utils"
	"github.com/iwinardhyas/sprint_backend/platform/database"
)

func Cat(c *fiber.Ctx) error {
	cats := &models.Cat{}

	// Checking received data from JSON body.
	if err := c.BodyParser(cats); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(cats); err != nil {
		// Return, if some fields are not valid.
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	cat := &models.Cat{}

	cat.ID = uuid.New()
	cat.Name = cats.Name
	cat.Race = cats.Race
	cat.Sex = cats.Sex
	cat.AgeInMonth = cats.AgeInMonth
	cat.Description = cats.Description
	cat.ImageUrls = cats.ImageUrls

	// Validate user fields.
	if err := validate.Struct(cat); err != nil {
		// Return, if some fields are not valid.
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := db.CreateCat(cat); err != nil {
		// Return status 500 and create user process error.
		return c.Status(409).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"email": "oke",
		},
	})
}
