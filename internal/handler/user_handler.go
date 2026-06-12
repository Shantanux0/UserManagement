package handler

import (
	"errors"
	"strconv"
	"user-management/internal/models"
	"user-management/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON body",
		})
	}

	if err := Validate.Struct(req); err != nil {
		h.logger.Warn("Validation failed for CreateUser", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	h.logger.Info("Created user successfully", zap.Int32("id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid user ID param", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID parameter",
		})
	}

	user, err := h.service.GetUserByID(c.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			h.logger.Warn("User not found", zap.Int("id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		h.logger.Error("Failed to get user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid user ID param", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID parameter",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON body",
		})
	}

	if err := Validate.Struct(req); err != nil {
		h.logger.Warn("Validation failed for UpdateUser", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.service.UpdateUser(c.Context(), int32(id), req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			h.logger.Warn("User not found for update", zap.Int("id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		h.logger.Error("Failed to update user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	h.logger.Info("Updated user successfully", zap.Int32("id", user.ID))
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid user ID param", zap.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID parameter",
		})
	}

	err = h.service.DeleteUser(c.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			h.logger.Warn("User not found for delete", zap.Int("id", id))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		h.logger.Error("Failed to delete user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	h.logger.Info("Deleted user successfully", zap.Int("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var page, limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid page parameter",
			})
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid limit parameter",
			})
		}
	}

	users, err := h.service.ListUsers(c.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users",
		})
	}

	return c.JSON(users)
}
