package service

import (
	"avito-merch/internal/model"
	"avito-merch/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo    repository.UserRepository
	secret      string
	jwtLifetime time.Duration
}

func NewAuthService(userRepo repository.UserRepository, secret string) AuthService {
	return &authService{
		userRepo:    userRepo,
		secret:      secret,
		jwtLifetime: 24 * time.Hour,
	}
}

func (s *authService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return s.createNewUser(ctx, username, password)
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.GenerateToken(user.ID)
}

func (s *authService) createNewUser(ctx context.Context, username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Coins:        model.InitialCoins,
	}

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	if err := s.userRepo.Create(ctx, tx, user); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return s.GenerateToken(user.ID)
}

func (s *authService) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.jwtLifetime).Unix(),
	})

	return token.SignedString([]byte(s.secret))
}

func (s *authService) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return 0, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["sub"].(float64)), nil
	}

	return 0, ErrInvalidToken
}
