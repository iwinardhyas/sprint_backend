package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iwinardhyas/sprint_backend/app/controllers"
	"github.com/iwinardhyas/sprint_backend/pkg/middleware"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)   // renew Access & Refresh tokens

}
