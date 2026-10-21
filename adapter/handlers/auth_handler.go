package handlers

import (
	"errors"
	"log"
	"time"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/config"
	"expense.com/m/entity"
	"expense.com/m/module/auth"
	"github.com/gofiber/fiber/v2"
)

const refreshCookieName = "refresh_token"

type AuthHandler struct {
	service *auth.Service
	cfg     *config.AuthConfig
}

func NewAuthHandler(service *auth.Service, cfg *config.AuthConfig) *AuthHandler {
	return &AuthHandler{service: service, cfg: cfg}
}

type authResponse struct {
	AccessToken string                   `json:"accessToken"`
	ExpiresIn   int64                    `json:"expiresIn"`
	User        *presenter.UserPresenter `json:"user"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	result, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return base_response.Error(c, fiber.StatusUnauthorized, "Email atau kata sandi salah")
		}
		log.Printf("login error: %v", err)
		return base_response.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server")
	}

	h.setRefreshCookie(c, result.RefreshToken)

	return base_response.Success(c, toAuthResponse(result))
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	result, err := h.service.Refresh(c.Context(), c.Cookies(refreshCookieName))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			h.clearRefreshCookie(c)
			return base_response.Error(c, fiber.StatusUnauthorized, "Sesi tidak valid, silakan login kembali")
		}
		log.Printf("refresh error: %v", err)
		return base_response.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server")
	}

	h.setRefreshCookie(c, result.RefreshToken)

	return base_response.Success(c, toAuthResponse(result))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.service.Logout(c.Context(), c.Cookies(refreshCookieName)); err != nil {
		log.Printf("logout error: %v", err)
	}

	h.clearRefreshCookie(c)

	return base_response.Success(c, fiber.Map{"message": "Logout berhasil"})
}

func (h *AuthHandler) setRefreshCookie(c *fiber.Ctx, value string) {
	if value == "" {
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: h.cfg.CookieSameSite,
		MaxAge:   int(h.cfg.RefreshTTL.Seconds()),
		Expires:  time.Now().Add(h.cfg.RefreshTTL),
	})
}

func (h *AuthHandler) clearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: h.cfg.CookieSameSite,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
	})
}

func toAuthResponse(result *auth.AuthResult) authResponse {
	return authResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
		User:        new(presenter.UserPresenter).FromDTO(result.User),
	}
}
