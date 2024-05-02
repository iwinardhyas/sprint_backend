package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewCat(c *fiber.Ctx) error {
	newCat := &models.NewCat{}
	if err := c.BodyParser(newCat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	validate := utils.NewValidator()

	if err := validate.Struct(newCat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userID := claims.UserID

	cat := &models.Cat{}
	cat.ID = uuid.New()
	cat.UserID = userID
	cat.Name = newCat.Name
	cat.Race = newCat.Race
	cat.Sex = newCat.Sex
	cat.AgeInMonth = newCat.AgeInMonth
	cat.ImageUrls = newCat.ImageUrls
	cat.Description = newCat.Description
	cat.HasMatched = false
	cat.CreatedAt = time.Now()

	if err := validate.Struct(cat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}
	if err := db.CreateCat(cat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"id":        cat.ID,
			"createdAt": cat.CreatedAt,
		},
	})
}

func GetCats(c *fiber.Ctx) error {
	query := ""

	id := c.Query("id")
	race := c.Query("race")
	sex := c.Query("sex")
	hasMatchedStr := c.Query("hasMatched")
	ageInMonthStr := c.Query("ageInMonth")
	ownedStr := c.Query("owned")
	search := c.Query("search")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if id != "" {
		query += fmt.Sprintf(" AND id = '%s'", id)
	}

	if race != "" {
		allowedRaces := map[string]bool{
			"persian":           true,
			"maine coon":        true,
			"siamese":           true,
			"ragdoll":           true,
			"bengal":            true,
			"sphynx":            true,
			"british shorthair": true,
			"abyssinian":        true,
			"scottish fold":     true,
			"birman":            true,
		}
		if allowedRaces[strings.ToLower(race)] {
			query += fmt.Sprintf(" AND race ILIKE '%s'", race)
		}
	}

	if sex != "" {
		sex_lower := strings.ToLower(sex)
		if sex_lower == "male" || sex_lower == "female" {
			query += fmt.Sprintf(" AND sex = '%s'", sex_lower)
		}
	}

	if hasMatchedStr == "true" || hasMatchedStr == "false" {
		hasMatched, _ := strconv.ParseBool(hasMatchedStr)
		query += fmt.Sprintf(" AND has_matched = %t", hasMatched)
	}

	if ageInMonthStr != "" {
		var comparison = "="
		var value int
		var err error
		if ageInMonthStr[0] == '<' {
			value, err = strconv.Atoi(ageInMonthStr[1:])
			comparison = "<"
		} else if ageInMonthStr[0] == '>' {
			value, err = strconv.Atoi(ageInMonthStr[1:])
			comparison = ">"
		} else {
			value, err = strconv.Atoi(ageInMonthStr)
		}
		if err == nil {
			query += fmt.Sprintf(" AND age_in_month %s %d", comparison, value)
		}
	}

	if ownedStr == "true" || ownedStr == "false" {
		claims, err := utils.ExtractTokenMetadata(c)
		if err == nil {
			userID := claims.UserID.String()
			query += fmt.Sprintf(" AND user_id = %s", userID)
		}
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query += fmt.Sprintf(` AND name ILIKE '%[1]s' OR race ILIKE '%[1]s' OR description ILIKE '%[1]s'`, searchTerm)
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			query += fmt.Sprintf(" LIMIT %d", limit)
		}
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	db, err := database.OpenDBConnection()

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	res, err := db.GetCats(query)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    res,
	})
}

func DeleteCat(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userID := claims.UserID.String()

	catID := c.Params("id")

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	err = db.DeleteCat(catID, userID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error()})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Deletion successfull",
	})
}
