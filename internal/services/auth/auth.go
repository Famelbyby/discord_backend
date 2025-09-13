package auth

import (
	"context"
	"errors"
	"fmt"
	"idea-store-auth/internal/domain/models"
	"idea-store-auth/internal/lib/jwt"
	"idea-store-auth/internal/lib/logger/sl"
	"idea-store-auth/internal/storage"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passwordHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// New returns a new instance of the Auth service.
func New(log *slog.Logger, userProvider UserProvider, userSaver UserSaver, tokenTTL time.Duration) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("servic Login error: " + ErrInvalidCredentials.Error())
		}

		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("servic Login error: " + ErrInvalidCredentials.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("servic Login error: " + ErrInvalidCredentials.Error())
	}

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("Failed to generate token", sl.Err(err))
		return "", fmt.Errorf("servic Login error: " + err.Error())
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	slog.Info("Registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		slog.Error("Failed to generate password hash", sl.Err(err))
		return -1, fmt.Errorf("servic RegisterNewUser error: " + err.Error())
	}
	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			slog.Warn("User already exists")

			return -1, fmt.Errorf("servic RegisterNewUser error: " + storage.ErrUserExists.Error())
		}
		slog.Error("Failed to save user", sl.Err(err))
		return -1, fmt.Errorf("servic Login error: " + err.Error())
	}

	return id, nil
}
