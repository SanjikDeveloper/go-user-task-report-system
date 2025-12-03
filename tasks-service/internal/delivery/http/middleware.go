package delivery

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}

func (h *Handler) CheckJWT(c *gin.Context) {
	h.logger.Info("CheckJWT: starting authentication check")

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		h.logger.Error("CheckJWT: Authorization header is missing")
		c.AbortWithStatus(401)
		return
	}
	h.logger.Info("CheckJWT: Authorization header received, length: %d", len(header))

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		h.logger.Error("CheckJWT: Authorization header format is invalid, parts count: %d", len(headerParts))
		c.AbortWithStatus(401)
		return
	}

	if len(headerParts[1]) == 0 {
		h.logger.Error("CheckJWT: Token string is empty")
		c.AbortWithStatus(401)
		return
	}

	tokenString := headerParts[1]
	h.logger.Info("CheckJWT: Token extracted, length: %d, first 20 chars: %s", len(tokenString), tokenString[:min(20, len(tokenString))])

	userId, err := h.parseToken(tokenString)
	if err != nil {
		h.logger.Error("CheckJWT: parseToken failed: %v", err)
		c.AbortWithStatus(401)
		return
	}

	h.logger.Info("CheckJWT: Token parsed successfully, userId: %d", userId)
	c.Set(userCtx, userId)
	h.logger.Info("CheckJWT: userId set in context")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *Handler) parseToken(tokenString string) (int, error) {
	h.logger.Info("parseToken: starting token parsing")

	signingKey := h.services.GetJWTSigningKey()
	if signingKey == "" {
		h.logger.Error("parseToken: JWT signing key is empty!")
		return 0, errors.New("JWT signing key is not configured")
	}
	h.logger.Info("parseToken: JWT signing key retrieved, length: %d", len(signingKey))

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		h.logger.Info("parseToken: checking signing method, alg: %v", token.Method.Alg())
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.logger.Error("parseToken: invalid signing method: %v", token.Method)
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		h.logger.Error("parseToken: jwt.ParseWithClaims failed: %v", err)
		return 0, errors.Wrap(err, "jwt.ParseWithClaims() err:")
	}

	if !token.Valid {
		h.logger.Error("parseToken: token is not valid")
		return 0, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		h.logger.Error("parseToken: failed to cast claims to *tokenClaims, type: %T", token.Claims)
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	h.logger.Info("parseToken: claims extracted successfully, userId: %d", claims.UserId)
	return claims.UserId, nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.Errorf("user id has wrong type: %T, value: %v", id, id)
	}

	return idInt, nil
}
