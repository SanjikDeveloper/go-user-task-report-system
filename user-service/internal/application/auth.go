package application

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"user-service/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	tokenTTL = 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}

func (s *Service) CreateUser(ctx context.Context, user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password, s.auth.PasswordSalt)
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, username, generatePasswordHash(password, s.auth.PasswordSalt))
	if err != nil {

		if errors.Is(err, models.ErrNotFound) {
			return "", models.ErrUnauthorized
		}
		return "", errors.Wrap(err, "s.post.GetUser() err:")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(s.auth.JWTSigningKey))
}

func (s *Service) ParseToken(acessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(acessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.auth.JWTSigningKey), nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "jwt.ParseWithClaims() err:")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
