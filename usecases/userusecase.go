package usecases

import (
	"errors"
	"fmt"
	"loan/domain"
	"loan/infrastructure"
	"loan/repositories"
	"loan/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase interface {
	Register(req domain.RegisterRequest) error
	VerifyEmail(email, token string) error
	Login(*domain.AuthUser) (string, string, error)
	GetAllUsers() ([]*domain.User, error)
	DeleteUser(objectID primitive.ObjectID) error
}

type userUsecase struct {
	userRepo repositories.UserRepository
	jwtSvc   infrastructure.JWTService
}

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {
	return &userUsecase{userRepo: repo, jwtSvc: jr}
}

func (u *userUsecase) Register(req domain.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:      req.Email,
		Password:   hashedPassword,
		Profile:    req.Profile,
		IsVerified: false,
	}

	// Create user in MongoDB
	err = u.userRepo.Create(user)
	if err != nil {
		return err
	}

	// Generate email verification token
	token, err := utils.GenerateEmailToken(user.Email)
	if err != nil {
		return err
	}

	// Send verification email
	err = utils.SendVerificationEmail(user.Email, token)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) VerifyEmail(email, token string) error {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("invalid email")
	}

	// Check if token is valid
	if !utils.ValidateEmailToken(token, user.Email) {
		return errors.New("invalid token")
	}

	// Mark the user's account as verified
	user.IsVerified = true
	return u.userRepo.Update(user)
}

// Login authenticates a user and returns JWT and refresh tokens if successful
func (u *userUsecase) Login(authUser *domain.AuthUser) (string, string, error) {
	fmt.Println("authuser: ", authUser)
	user, err := u.userRepo.GetUserByUsername(authUser.Username)
	if err != nil {
		return "", "", err
	}

	fmt.Println("user: ", user)

	if err := utils.CheckPasswordHash(user.Password, authUser.Password); err != nil {
		return "", "", errors.New("invalid username or password2")
	}

	// Generate JWT and refresh tokens for the authenticated user
	token, err := u.jwtSvc.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.jwtSvc.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshedTokenClaim := &domain.RefreshToken{
		UserID:    user.ID,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	// Save the refresh token in the database
	err = u.tokenRepo.SaveRefreshToken(refreshedTokenClaim)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	users, err := u.userRepo.GetAllUsers()
	return users, err
}

// DeleteUser deletes a user by ID
func (u *UserUsecase) DeleteUser(objectID primitive.ObjectID) error {
	return u.userRepo.DeleteUser(objectID)
}
