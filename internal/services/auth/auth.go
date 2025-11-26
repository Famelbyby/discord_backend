package auth

import (
	"context"
	"discord_backend/internal/domain/models"
	"discord_backend/internal/lib/logger/sl"
	"discord_backend/internal/storage"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxAge = time.Hour * 24 * 7
)

type Auth struct {
	log             *slog.Logger
	userSaver       UserSaver
	userProvider    UserProvider
	tokenTTL        time.Duration
	sessionProvider SessionProvider
}

type SessionProvider interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, key ...string) *redis.IntCmd
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
func New(log *slog.Logger, userProvider UserProvider, userSaver UserSaver, sessionProvider SessionProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		userSaver:       userSaver,
		userProvider:    userProvider,
		log:             log,
		tokenTTL:        tokenTTL,
		sessionProvider: sessionProvider,
	}
}

func AddSession(a *Auth, ctx context.Context, id string) (string, error) {
	existedSessionId, err := a.sessionProvider.Get(ctx, id).Result()

	if err != redis.Nil {
		_, err := a.sessionProvider.Del(ctx, id).Result()

		if err != nil {
			return "", fmt.Errorf("service Login error: " + err.Error())
		}

		_, err = a.sessionProvider.Del(ctx, existedSessionId).Result()

		if err != nil {
			return "", fmt.Errorf("service Login error: " + err.Error())
		}
	}

	var sessionId string

	for {
		sessionId = uuid.NewString()

		_, err := a.sessionProvider.Get(ctx, sessionId).Result()

		if err == redis.Nil {
			break
		}

		if err != nil {
			return "", fmt.Errorf("service redis error:" + err.Error())
		}
	}

	err = a.sessionProvider.Set(ctx, sessionId, id, maxAge).Err()

	if err != nil {
		a.log.Error("failed to set session", sl.Err(err))

		return "", fmt.Errorf("service redis error: " + err.Error())
	}

	err = a.sessionProvider.Set(ctx, id, sessionId, maxAge).Err()

	if err != nil {
		a.log.Error("failed to set session", sl.Err(err))

		return "", fmt.Errorf("service redis error: " + err.Error())
	}

	return sessionId, nil
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	id string,
) (string, error) {
	slog.Info("logining")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("service Login error: " + ErrInvalidCredentials.Error())
		}

		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("service Login error: " + ErrInvalidCredentials.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("service Login error: " + ErrInvalidCredentials.Error())
	}

	existedSessionId, err := a.sessionProvider.Get(ctx, id).Result()

	if err != redis.Nil {
		_, err := a.sessionProvider.Del(ctx, id).Result()

		if err != nil {
			return "", fmt.Errorf("service Login error: " + err.Error())
		}

		_, err = a.sessionProvider.Del(ctx, existedSessionId).Result()

		if err != nil {
			return "", fmt.Errorf("service Login error: " + err.Error())
		}
	}

	sessionId, err := AddSession(a, ctx, id)

	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
	id string,
) (string, error) {
	slog.Info("Registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		slog.Error("Failed to generate password hash", sl.Err(err))
		return "", fmt.Errorf("servic RegisterNewUser error: " + err.Error())
	}
	_, err = a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			slog.Warn("User already exists")

			return "", fmt.Errorf("servic RegisterNewUser error: " + storage.ErrUserExists.Error())
		}
		slog.Error("Failed to save user", sl.Err(err))
		return "", fmt.Errorf("servic Login error: " + err.Error())
	}

	sessionId, err := AddSession(a, ctx, id)

	fmt.Println(sessionId)

	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (a *Auth) Logout(
	ctx context.Context,
	sessionId string,
) error {
	slog.Info("logouting")

	userId, err := a.sessionProvider.Get(ctx, sessionId).Result()

	if err == redis.Nil {
		slog.Error("failed to logout", sl.Err(err))

		return fmt.Errorf("service login error: no session")
	}

	err = a.sessionProvider.Del(ctx, sessionId).Err()

	if err != nil {
		slog.Error("failed to logout", sl.Err(err))

		return fmt.Errorf("service login error: " + err.Error())
	}

	err = a.sessionProvider.Del(ctx, userId).Err()

	if err != nil {
		slog.Error("failed to logout", sl.Err(err))

		return fmt.Errorf("service login error: " + err.Error())
	}

	return nil
}

func (a *Auth) IsRegistered(
	ctx context.Context,
	sessionId string,
) (string, error) {
	slog.Info("logouting")

	userId, err := a.sessionProvider.Get(ctx, sessionId).Result()

	if err == redis.Nil {
		slog.Error("failed to logout", sl.Err(err))

		return "", fmt.Errorf("service login error: no session")
	}

	return userId, nil
}
