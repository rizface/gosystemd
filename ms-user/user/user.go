package user

import (
	"context"
	"errors"
	"os"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Inserter interface {
	Insert(context.Context, User) (User, error)
}

type Getter interface {
	GetByUsername(context.Context, string) (User, error)
}

type UserRepo interface {
	Inserter
	Getter
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Claim struct {
	jwt.RegisteredClaims
	UserId string `json:"userId"`
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

func (u *User) GenerateID() {
	u.Id = uuid.NewString()
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Username, validation.Required),
	)
}

func (u User) GetJWt() (string, error) {
	accessTokenSecret := os.Getenv("SYSTEMD_JWT_ACCESS_SECRET")
	if accessTokenSecret == "" {
		return "", errors.New("empty jwt access secrett")
	}

	claim := Claim{
		UserId: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rizface",
			Subject:   "auth token",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(accessTokenSecret))
}

func (u User) VerifyPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateUser(ctx context.Context, repo UserRepo, newUser User) (User, error) {
	if err := newUser.Validate(); err != nil {
		return User{}, err
	}

	user, err := repo.GetByUsername(ctx, newUser.Username)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return User{}, err
	}

	if user.Username != "" {
		return User{}, ErrUsernameAlreadyTaken
	}

	newUser.HashPassword()
	newUser.GenerateID()

	user, err = repo.Insert(ctx, newUser)
	if err != nil {
		return User{}, err
	}

	user.Password = ""

	return user, nil
}

func Login(ctx context.Context, repo UserRepo, usr User) (string, error) {
	user, err := repo.GetByUsername(ctx, usr.Username)
	if err != nil {
		return "", err
	}

	match, err := user.VerifyPassword(usr.Password)
	if err != nil {
		return "", err
	}

	if !match {
		return "", ErrWrongPassword
	}

	// [TODO]: Store token to redis and expired existing session if user relogin

	return user.GetJWt()
}
