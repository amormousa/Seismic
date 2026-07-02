package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/helpers"
	"github.com/majoramari/seismic/apps/api/models"
	"github.com/majoramari/seismic/apps/api/services"
)

type AuthHandler struct {
	Pool           *pgxpool.Pool
	EmailCfg       services.EmailConfig
	JWTSecret      string
	JWTExpiryHours int
}

type magicLinkRequest struct {
	Email string `json:"email"`
}

// RequestMagicLink handles POST /api/auth/magic-link.
// Creates a login token tied to the email and sends it,
// regardless of whether the user already exists.
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

// VerifyMagicLink handles GET /api/auth/verify.
// Validates the token. If the user already exists, logs them
// in. If not, returns a signup token so the frontend can show
// an onboarding screen to pick a username.
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

	jwtToken, err := generateJWT(user.ID, h.JWTSecret, h.JWTExpiryHours)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	return helpers.Success(c, "Logged in successfully", fiber.Map{
		"newUser": false,
		"token":   jwtToken,
		"user":    user,
	})
}

type completeSignupRequest struct {
	SignupToken string `json:"signupToken"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

// CompleteSignup handles POST /api/auth/complete-signup.
// Finishes creating the account once a new user picks a
// username, using the signup token from VerifyMagicLink.
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

	jwtToken, err := generateJWT(user.ID, h.JWTSecret, h.JWTExpiryHours)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	return helpers.Success(c, "Account created", fiber.Map{
		"token": jwtToken,
		"user":  user,
	})
}
