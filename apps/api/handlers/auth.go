package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/helpers"
	"github.com/majoramari/seismic/apps/api/models"
	"github.com/majoramari/seismic/apps/api/services"
)

type AuthHandler struct {
	Pool      *pgxpool.Pool
	EmailCfg  services.EmailConfig
	JWTSecret string
}

type magicLinkRequest struct {
	Email string `json:"email"`
}

// RequestMagicLink godoc
// @Summary      Request a magic link
// @Description  Sends a login link to the given email. Works for both new and existing users.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body magicLinkRequest true "Email address"
// @Success      200 {object} helpers.APIResponse
// @Failure      400 {object} helpers.APIResponse
// @Router       /api/auth/magic-link [post]
func (h *AuthHandler) RequestMagicLink(c *fiber.Ctx) error {
	var body magicLinkRequest
	if err := c.BodyParser(&body); err != nil {
		return helpers.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	email := strings.TrimSpace(strings.ToLower(body.Email))
	if email == "" || !strings.Contains(email, "@") {
		return helpers.Error(c, fiber.StatusBadRequest, "Please provide a valid email address")
	}

	ctx := c.Context()

	link, err := models.CreateMagicLink(ctx, h.Pool, email)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create login link")
	}

	err = services.SendMagicLinkEmail(h.EmailCfg, email, link.Token)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to send login email")
	}

	return helpers.Success(c, "Check your email for a login link", nil)
}

// VerifyMagicLink godoc
// @Summary      Verify a magic link token
// @Description  Validates the token from the login email. Logs in existing users, or returns a signup token for new users.
// @Tags         auth
// @Produce      json
// @Param        token query string true "Magic link token"
// @Success      200 {object} helpers.APIResponse
// @Failure      400 {object} helpers.APIResponse
// @Failure      401 {object} helpers.APIResponse
// @Router       /api/auth/verify [get]
func (h *AuthHandler) VerifyMagicLink(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return helpers.Error(c, fiber.StatusBadRequest, "Missing token")
	}

	ctx := c.Context()

	link, err := models.FindMagicLinkByToken(ctx, h.Pool, token)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}
	if link == nil {
		return helpers.Error(c, fiber.StatusUnauthorized, "Invalid login link")
	}
	if link.Used {
		return helpers.Error(c, fiber.StatusUnauthorized, "This login link has already been used")
	}
	if link.IsExpired() {
		return helpers.Error(c, fiber.StatusUnauthorized, "This login link has expired")
	}

	if err := models.MarkMagicLinkUsed(ctx, h.Pool, link.ID); err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}

	user, err := models.FindUserByEmail(ctx, h.Pool, link.Email)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}

	// New user, needs to pick a username before account is created
	if user == nil {
		signupToken, err := generateSignupToken(link.Email, h.JWTSecret)
		if err != nil {
			return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
		}
		return helpers.Success(c, "New account, complete signup", fiber.Map{
			"newUser":     true,
			"signupToken": signupToken,
			"email":       link.Email,
		})
	}

	// Existing user, log in normally
	accessToken, err := generateJWT(user.ID, h.JWTSecret)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	refreshToken, err := models.CreateRefreshToken(ctx, h.Pool, user.ID)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create session")
	}
	setRefreshTokenCookie(c, refreshToken)

	return helpers.Success(c, "Logged in successfully", fiber.Map{
		"newUser":     false,
		"accessToken": accessToken,
		"user":        user,
	})
}

type completeSignupRequest struct {
	SignupToken string `json:"signupToken"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

// CompleteSignup godoc
// @Summary      Complete signup for a new user
// @Description  Creates the account using the signup token from VerifyMagicLink, once a username is chosen.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body completeSignupRequest true "Signup details"
// @Success      200 {object} helpers.APIResponse
// @Failure      400 {object} helpers.APIResponse
// @Failure      401 {object} helpers.APIResponse
// @Failure      409 {object} helpers.APIResponse
// @Router       /api/auth/complete-signup [post]
func (h *AuthHandler) CompleteSignup(c *fiber.Ctx) error {
	var body completeSignupRequest
	if err := c.BodyParser(&body); err != nil {
		return helpers.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	email, err := verifySignupToken(body.SignupToken, h.JWTSecret)
	if err != nil {
		return helpers.Error(c, fiber.StatusUnauthorized, "Invalid or expired signup token")
	}

	username := strings.TrimSpace(strings.ToLower(body.Username))
	if username == "" {
		return helpers.Error(c, fiber.StatusBadRequest, "Username is required")
	}

	ctx := c.Context()

	existing, _ := models.FindUserByUsername(ctx, h.Pool, username)
	if existing != nil {
		return helpers.Error(c, fiber.StatusConflict, "Username is already taken")
	}

	user, err := models.CreateUser(ctx, h.Pool, email, username, body.DisplayName)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create account")
	}

	accessToken, err := generateJWT(user.ID, h.JWTSecret)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	refreshToken, err := models.CreateRefreshToken(ctx, h.Pool, user.ID)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create session")
	}
	setRefreshTokenCookie(c, refreshToken)

	return helpers.Success(c, "Account created", fiber.Map{
		"accessToken": accessToken,
		"user":        user,
	})
}

// RefreshAccessToken godoc
// @Summary      Refresh access token
// @Description  Uses the refresh token cookie to issue a new access token and rotate the refresh token.
// @Tags         auth
// @Produce      json
// @Success      200 {object} helpers.APIResponse
// @Failure      401 {object} helpers.APIResponse
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	rawToken := c.Cookies("refresh_token")
	if rawToken == "" {
		return helpers.Error(c, fiber.StatusUnauthorized, "No refresh token provided")
	}

	ctx := c.Context()

	stored, err := models.FindValidRefreshToken(ctx, h.Pool, rawToken)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}
	if stored == nil {
		return helpers.Error(c, fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}

	if err := models.RevokeRefreshToken(ctx, h.Pool, stored.ID); err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}

	newRawToken, err := models.CreateRefreshToken(ctx, h.Pool, stored.UserID)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}
	setRefreshTokenCookie(c, newRawToken)

	accessToken, err := generateJWT(stored.UserID, h.JWTSecret)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
	}

	return helpers.Success(c, "Token refreshed", fiber.Map{
		"accessToken": accessToken,
	})
}

// setRefreshTokenCookie sets the refresh token as an httpOnly
// cookie so client-side JavaScript can never read it directly.
func setRefreshTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/api/auth",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
}

// GetAPIKey godoc
// @Summary      Get API key
// @Description  Returns the logged-in user's API key, used to configure editor plugins.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} helpers.APIResponse
// @Failure      401 {object} helpers.APIResponse
// @Router       /api/auth/apikey [get]
func (h *AuthHandler) GetAPIKey(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx := c.Context()
	user, err := models.FindUserByID(ctx, h.Pool, userID)
	if err != nil || user == nil {
		return helpers.Error(c, fiber.StatusNotFound, "User not found")
	}

	return helpers.Success(c, "API key retrieved", fiber.Map{
		"apiKey": user.APIKey,
	})
}

// RegenerateAPIKey godoc
// @Summary      Regenerate API key
// @Description  Replaces the user's API key, invalidating any editor plugins still using the old one.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} helpers.APIResponse
// @Failure      401 {object} helpers.APIResponse
// @Router       /api/auth/apikey/regenerate [post]
func (h *AuthHandler) RegenerateAPIKey(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx := c.Context()
	newKey, err := models.RegenerateAPIKey(ctx, h.Pool, userID)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to regenerate api key")
	}

	return helpers.Success(c, "API key regenerated", fiber.Map{
		"apiKey": newKey,
	})
}
