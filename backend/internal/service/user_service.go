package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(email, password string) (*domain.User, error)
	Login(email, password string) (string, string, error)
	Refresh(refreshToken string) (string, string, error)
	GetUserByID(id uint) (*domain.User, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
}

type userService struct {
	repo      repository.UserRepository
	resetRepo repository.PasswordResetRepository
}

func NewUserService(repo repository.UserRepository, resetRepo repository.PasswordResetRepository) UserService {
	return &userService{repo: repo, resetRepo: resetRepo}
}

func (s *userService) Register(email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func generateTokens(userID uint) (string, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	aToken, err := accessToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"refresh": true,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	rToken, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return aToken, rToken, nil
}

func (s *userService) Login(email, password string) (string, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	return generateTokens(user.ID)
}

func (s *userService) Refresh(refreshToken string) (string, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// Check if it's a refresh token
	if _, ok := claims["refresh"]; !ok {
		return "", "", errors.New("not a refresh token")
	}

	userID := uint(claims["user_id"].(float64))
	return generateTokens(userID)
}

func (s *userService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) ForgotPassword(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil // Avoid exposing if email exists
	}

	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	tokenStr := hex.EncodeToString(bytes)

	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.resetRepo.CreateToken(resetToken); err != nil {
		return err
	}

	// Simulate email by logging
	log.Printf("[SIMULATED EMAIL] Password reset requested for %s. Token: %s", email, tokenStr)
	return nil
}

func (s *userService) ResetPassword(tokenStr, newPassword string) error {
	token, err := s.resetRepo.FindToken(tokenStr)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(token.ExpiresAt) {
		s.resetRepo.DeleteToken(token.ID)
		return errors.New("token expired")
	}

	user, err := s.repo.FindByID(token.UserID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Invalidate token
	s.resetRepo.DeleteToken(token.ID)
	return nil
}
