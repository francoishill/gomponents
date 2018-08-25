package token

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/francoishill/gomponents/user"
)

type Service interface {
	Middlewares() []func(http.Handler) http.Handler

	Create(user user.User) (string, error)
	UserIDFromContext(ctx context.Context) (string, error)
}

func JWTService(signKey []byte, expiryDuration time.Duration, addUserInfoToClaimsFunc func(claims jwtauth.Claims, user user.User) error) *jwtService {
	if len(signKey) == 0 {
		logrus.Panic("signKey is required in InitTokenAuth")
	}
	auth := jwtauth.New("HS256", signKey, nil)

	return &jwtService{
		auth,
		expiryDuration,
		addUserInfoToClaimsFunc,
	}
}

type jwtService struct {
	auth                    *jwtauth.JWTAuth
	expiryDuration          time.Duration
	addUserInfoToClaimsFunc func(claims jwtauth.Claims, user user.User) error
}

func (t *jwtService) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		jwtauth.Verifier(t.auth), // Seek, verify and validate JWT tokens (only sets invalid token error on context but continues, the Authenticator errors on invalid token)
		jwtauth.Authenticator,    // Handle valid / invalid tokens
	}
}

func (t *jwtService) Create(user user.User) (string, error) {
	//refer to github.com/dgrijalva/jwt-go->StandardClaims and https://tools.ietf.org/html/rfc7519#section-4.1
	claims := jwtauth.Claims{
		"iat": time.Now().Unix(),                       //IssuedAt
		"exp": time.Now().Add(t.expiryDuration).Unix(), //ExpiresAt
	}

	if err := t.addUserInfoToClaimsFunc(claims, user); err != nil {
		return "", errors.Wrapf(err, "Failed to add user info to token claims")
	}

	_, tokenString, err := t.auth.Encode(claims)
	return tokenString, err
}

func (t *jwtService) UserIDFromContext(ctx context.Context) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get token from request context")
	}

	tmpUserID, ok := claims["user_id"]
	if !ok {
		userMessage := "Invalid token, user_id is missing"
		return "", errors.Errorf(userMessage)
	}

	userID, isStr := tmpUserID.(string)
	if !isStr {
		userMessage := "Invalid token, user_id is found but not of type string"
		return "", errors.Errorf(userMessage)
	}

	return userID, nil
}
