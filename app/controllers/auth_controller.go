package controllers

import (
	"context"
	"time"

	"github.com/iwinardhyas/sprint_backend/app/models"
	"github.com/iwinardhyas/sprint_backend/pkg/utils"
	"github.com/iwinardhyas/sprint_backend/platform/cache"
	"github.com/iwinardhyas/sprint_backend/platform/database"

	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserSignUp method to create a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param user_role body string true "User role"
// @Success 200 {object} models.User
// @Router /v1/user/sign/up [post]
func UserSignUp(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signUp := &models.SignUp{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signUp); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(signUp); err != nil {
		// Return, if some fields are not valid.
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	hash := sha256.New()
	// Create a new now date and time string with salt.
	refresh := os.Getenv("JWT_REFRESH_KEY") + time.Now().String()
	// See: https://pkg.go.dev/io#Writer.Write
	_, err = hash.Write([]byte(refresh))
	if err != nil {
		// Return error, it refresh token generation failed.
		return err
	}
	// Set expires hours count for refresh key from .env file.
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
	// Set expiration time.
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	// Create a new refresh token (sha256 string with salt + expire time).
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime

	// Create a new user struct.
	user := &models.User{}

	// Set initialized default data for user:
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.Email = signUp.Email
	user.Name = signUp.Name
	user.PasswordHash = utils.GeneratePassword(signUp.Password)
	user.UserStatus = 1 // 0 == blocked, 1 == active
	user.UserRole = "user"
	user.AccessToken = t

	// Validate user fields.
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create a new user with validated data.
	if err := db.CreateUser(user); err != nil {
		// Return status 500 and create user process error.
		return c.Status(409).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Delete password hash field from JSON view.
	user.PasswordHash = ""

	// Return status 200 OK.
	return c.Status(201).JSON(fiber.Map{
		// "error": false,
		"message": "User registered successfully",
		"data":    user,
	})

}

// UserSignIn method to auth user and return access and refresh tokens.
// @Description Auth user and return access and refresh token.
// @Summary auth user and return access and refresh token
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "User Email"
// @Param password body string true "User Password"
// @Success 200 {string} status "ok"
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signIn := &models.SignIn{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(signIn); err != nil {
		// Return, if some fields are not valid.
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get user by email.
	foundedUser, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		// Return, if user not found.
		return c.Status(404).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given email is not found",
		})
	}

	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(foundedUser.PasswordHash, signIn.Password)
	if !compareUserPassword {
		// Return, if password is not compare to stored in database.
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   "wrong user email address or password",
		})
	}

	// Get role credentials from founded user.
	credentials, err := utils.GetCredentialsByRole(foundedUser.UserRole)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(foundedUser.ID.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.Status(201).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"email": signIn.Email,
			"token": tokens.Access,
			"name":  foundedUser.Name,
		},
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Description De-authorize user and delete refresh token from Redis.
// @Summary de-authorize user and delete refresh token from Redis
// @Tags User
// @Accept json
// @Produce json
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/user/sign/out [post]
func UserSignOut(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Save refresh token to Redis.
	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		// Return status 500 and Redis deletion error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errDelFromRedis.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
