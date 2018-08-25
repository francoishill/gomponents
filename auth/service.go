package auth

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/francoishill/gomponents/clienterror"
	"github.com/francoishill/gomponents/encryption"
	"github.com/francoishill/gomponents/rendering"
	"github.com/francoishill/gomponents/token"
	"github.com/francoishill/gomponents/user"
)

type Service interface {
	Register(user User) (token string, err error)
	Login(user User, password string) (token string, err error)
	MagicLogin(user User, magicToken string) (token string, err error)
}

func DefaultService(userFactory user.Factory, rendering rendering.Service, encryption encryption.Service, token token.Service) *defaultService {
	return &defaultService{
		userFactory,
		rendering,
		encryption,
		token,
	}
}

type defaultService struct {
	userFactory user.Factory
	rendering   rendering.Service
	encryption  encryption.Service
	token       token.Service
}

func (a *defaultService) Register(u User) (token string, err error) {
	logger := logrus.NewEntry(logrus.StandardLogger())

	userRepo := a.userFactory.Repo()
	if err = userRepo.Add(u); err != nil {
		if userRepo.IsDupErr(err) {
			userMessage := "Failed to add new user, user already exists"
			logger.WithError(err).Error(userMessage)
			return "", clienterror.NewErrorDefaultStatus(errors.New(userMessage))
		}
		userMessage := "Failed to add new user"
		logger.WithError(err).Error(userMessage)
		return "", errors.New(userMessage)
	}

	token, err = a.token.Create(u)
	if err != nil {
		userMessage := "Unable to generate token"
		logger.WithError(err).Error(userMessage)
		return "", errors.New(userMessage)
	}

	logger.Debug("Created token")
	return token, nil
}

func (a *defaultService) Login(user User, password string) (token string, err error) {
	logger := logrus.NewEntry(logrus.StandardLogger())

	if err := a.encryption.VerifyPassword(password, user.PasswordHash()); err != nil {
		logger.WithError(err).Error("User password mismatch")
		return "", errors.New("User email or password is incorrect")
	}
	logger = logger.WithField("user-id", user.ID())

	token, err = a.token.Create(user)
	if err != nil {
		userMessage := "Unable to generate token"
		logger.WithError(err).Error(userMessage)
		return "", errors.New(userMessage)
	}

	logger.Debug("Created token")
	return token, nil
}

func (a *defaultService) MagicLogin(user User, magicToken string) (token string, err error) {
	logger := logrus.NewEntry(logrus.StandardLogger()).WithField("user-id", user.ID())

	if user.MagicLoginToken() == nil {
		userMessage := fmt.Sprintf("Magic Token Login is not allowed for user with userID '%s' (since MagicLoginToken == NIL)", user.ID())
		logger.Error(userMessage)
		return "", errors.New(userMessage)
	}

	if *user.MagicLoginToken() != magicToken {
		userMessage := fmt.Sprintf("Token mismatch of user with userID '%s'", user.ID())
		logger.Error(userMessage)
		return "", errors.New(userMessage)
	}

	token, err = a.token.Create(user)
	if err != nil {
		userMessage := "Unable to generate token"
		logger.WithError(err).Error(userMessage)
		return "", errors.New(userMessage)
	}

	logger.Debug("Created token")
	return token, nil
}
