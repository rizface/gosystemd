package user

import (
	"context"
	"errors"

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
	Password string `json:"password"`
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

func CreateUser(ctx context.Context, repo UserRepo, newUser User) (User, error) {
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

	return user, nil
}
