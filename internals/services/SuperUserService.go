package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/htmltemplates"
	"github.com/lordofthemind/EventureGo/internals/newerrors"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/types"
	"github.com/lordofthemind/EventureGo/internals/utils"
	"github.com/lordofthemind/mygopher/gophersmtp"
	"github.com/lordofthemind/mygopher/gophertoken"
	"golang.org/x/crypto/bcrypt"
)

type SuperUserService struct {
	repo         repositories.SuperUserRepositoryInterface
	tokenManager gophertoken.TokenManager
	emailService gophersmtp.GopherSmtpInterface
}

func NewSuperUserService(
	repo repositories.SuperUserRepositoryInterface,
	tokenManager gophertoken.TokenManager,
	emailService gophersmtp.GopherSmtpInterface,
) SuperUserServiceInterface {
	return &SuperUserService{
		repo:         repo,
		tokenManager: tokenManager,
		emailService: emailService,
	}
}

// RegisterSuperUser registers a new superuser and sends verification email
func (s *SuperUserService) RegisterSuperUser(ctx context.Context, req *utils.RegisterSuperuserRequest) (*utils.RegisterSuperuserResponse, error) {
	// Validate email and username availability via validation layer
	if err := utils.ValidateUniqueness(ctx, req.Email, req.Username, s.repo); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to hash password")
	}

	role := "Guest"
	otp := uuid.New().String() // Generate a UUID as the OTP

	// Create a new SuperUser entity
	superUserEntity := types.NewSuperUser(req.Email, req.FullName, req.Username, string(hashedPassword), role, otp)

	// Store the SuperUser in the repository
	createdSuperUser, err := s.repo.CreateSuperUser(ctx, superUserEntity)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to create superuser")
	}

	// Generate verification link with OTP using the BaseURL from configuration
	verificationLink := fmt.Sprintf("%s/superusers/verify?otp=%s", configs.BaseURL, otp)

	// Prepare the email data
	emailData := map[string]interface{}{
		"FullName":         req.FullName,
		"VerificationLink": verificationLink,
		"OTP":              otp,
	}

	// Load and render the email template
	emailBody, err := htmltemplates.LoadAndRenderTemplate("welcome_verification_email.html", emailData)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to render email template")
	}

	// Send the verification email
	err = s.emailService.SendEmail([]string{req.Email}, "Welcome to EventureGo - Verify Your Account", emailBody, true)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to send verification email")
	}

	log.Printf("Generated OTP: %s", otp)

	// Prepare and return the response object
	return utils.CreateSuperuserResponse(createdSuperUser), nil
}

func (s *SuperUserService) VerifySuperUserOTP(ctx context.Context, otp string) (*types.SuperUserType, error) {
	// Fetch the superuser by OTP
	superUser, err := s.repo.FindSuperUserByOTP(ctx, otp)
	if err != nil {
		return nil, newerrors.Wrap(err, "invalid or expired OTP")
	}

	// Mark the superuser as verified
	if err := s.repo.VerifySuperUserOTP(ctx, superUser); err != nil {
		return nil, newerrors.Wrap(err, "failed to verify OTP")
	}

	// Send verification confirmation email
	confirmationEmailBody, err := htmltemplates.LoadAndRenderTemplate("verification_success_email.html", map[string]interface{}{
		"FullName": superUser.FullName,
	})
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to render confirmation email template")
	}

	err = s.emailService.SendEmail([]string{superUser.Email}, "Account Verified - Welcome to EventureGo", confirmationEmailBody, true)
	if err != nil {
		return nil, newerrors.Wrap(err, "failed to send verification confirmation email")
	}

	log.Printf("Superuser verified: %s", superUser.Email)

	return superUser, nil
}

func (s *SuperUserService) LogInSuperuser(ctx context.Context, loginRequest *utils.LogInSuperuserRequest) (*utils.LoginSuperuserResponse, error) {
	var superUser *types.SuperUserType
	var err error

	// Fetch the superuser by either email or username
	if loginRequest.Email != "" {
		superUser, err = s.repo.FindSuperUserByEmail(ctx, loginRequest.Email)
	} else if loginRequest.Username != "" {
		superUser, err = s.repo.FindSuperUserByUsername(ctx, loginRequest.Username)
	} else {
		return nil, newerrors.NewValidationError("email or username is required")
	}

	if err != nil {
		return nil, newerrors.NewValidationError("invalid email/username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(superUser.HashedPassword), []byte(loginRequest.Password)); err != nil {
		return nil, newerrors.NewValidationError("invalid email/username or password")
	}

	// Generate token with role
	authToken, err := s.tokenManager.GenerateToken(superUser.Username, configs.TokenExpiryDuration)
	if err != nil {
		return nil, err
	}

	// Prepare and return the response
	return &utils.LoginSuperuserResponse{
		ID:           superUser.ID,
		Email:        superUser.Email,
		Username:     superUser.Username,
		FullName:     superUser.FullName,
		Role:         superUser.Role,
		Token:        authToken,
		Is2FAEnabled: superUser.Is2FAEnabled,
	}, nil
}

func (s *SuperUserService) SendPasswordResetEmailWithUsernameOrEmail(ctx context.Context, email string, username string) error {
	var superUser *types.SuperUserType
	var err error

	// Fetch superuser by email or username
	if email != "" {
		superUser, err = s.repo.FindSuperUserByEmail(ctx, email)
	} else if username != "" {
		superUser, err = s.repo.FindSuperUserByUsername(ctx, username)
	} else {
		return newerrors.NewValidationError("email or username is required")
	}

	if err != nil {
		return newerrors.Wrap(err, "failed to find superuser")
	}

	// Generate and send a reset token
	resetToken := utils.GenerateResetToken()
	resetLink := fmt.Sprintf("%s/superuser/password-reset/%s", configs.BaseURL, resetToken)

	// Log the reset link for debugging
	log.Printf("Sending password reset token to %s: %s\n", superUser.Email, resetLink)

	// Prepare data for rendering email template
	templateData := map[string]interface{}{
		"FullName":            superUser.FullName,
		"ResetLink":           resetLink,
		"TokenExpiryDuration": configs.TokenExpiryDuration, // e.g., "15 minutes"
	}

	// Render the email template with the token data
	htmlBody, err := htmltemplates.LoadAndRenderTemplate("password_reset_email.html", templateData)
	if err != nil {
		return newerrors.Wrap(err, "failed to render email template")
	}

	// Send the password reset email using EmailService
	err = s.emailService.SendEmail([]string{superUser.Email}, "Password Reset Request", htmlBody, true)
	if err != nil {
		return newerrors.Wrap(err, "failed to send reset email")
	}

	// Store the reset token in the repository (assuming token expiry is also stored)
	return s.repo.UpdateResetToken(ctx, superUser.ID, resetToken)
}

// ResetPassword resets the password of a superuser using a token.
func (s *SuperUserService) ResetPassword(ctx context.Context, token, newPassword string) error {
	superUser, err := s.repo.FindSuperUserByResetToken(ctx, token)
	if err != nil {
		return newerrors.NewValidationError("invalid or expired reset token")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return newerrors.Wrap(err, "failed to hash new password")
	}

	// Update superuser's password and timestamp
	superUser.HashedPassword = string(hashedPassword)
	superUser.UpdatedAt = time.Now()
	superUser.ResetToken = nil

	// Update the superuser record
	return s.repo.UpdateSuperUser(ctx, superUser)
}

// // SeedSuperUser seeds a superuser based on the given request data
// func (s *SuperUserService) SeedSuperUser(ctx context.Context, req *utils.RegisterSuperuserRequest) error {
// 	// Check if the superuser already exists by email
// 	existingSuperUserByEmail, err := s.repo.FindSuperUserByEmail(ctx, req.Email)
// 	if err == nil && existingSuperUserByEmail != nil {
// 		return newerrors.NewValidationError("SuperUser with the provided email already exists")
// 	}

// 	// Check if the superuser already exists by username
// 	existingSuperUserByUsername, err := s.repo.FindSuperUserByUsername(ctx, req.Username)
// 	if err == nil && existingSuperUserByUsername != nil {
// 		return newerrors.NewValidationError("SuperUser with the provided username already exists")
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return newerrors.Wrap(err, "failed to hash password")
// 	}

// 	role := "SuperUser"
// 	// Create a new SuperUser entity
// 	superUserEntity := types.NewSuperUser(req.Email, req.FullName, req.Username, string(hashedPassword), role)

// 	// Store the SuperUser in the repository using the SuperUserService
// 	_, err = s.repo.CreateSuperUser(ctx, superUserEntity)
// 	if err != nil {
// 		return newerrors.Wrap(err, "failed to create superuser")
// 	}

// 	return nil
// }
