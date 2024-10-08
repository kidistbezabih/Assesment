package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"isactive"`
	IsAdmin   bool      `json:"isadmin"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type Token struct {
	ID           string `json:"id" bson:"_id"`
	UserId       string `json:"userid"`
	RefreshToken string `json:"refresh_token"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Email struct {
	User_email string `json:"email"`
}
type ResetForm struct {
	Passowrd    string `json:"password"`
	NewPassword string `json:"newpassword"`
}
type Profile struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type AuthRepository interface {
	CreateUser(ctx context.Context, user User) (string, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	DeleteUser(ctx context.Context, id string) error
	RegisterRefreshToken(ctx context.Context, userId string, token string) error
	GetRefreshToken(ctx context.Context, token string) (string, error)
	GetCollectionCount(ctx context.Context) (int64, error)
}

type AuthServices interface {
	Login(ctx context.Context, info LoginForm) (string, string, error)
	RegisterUser(ctx context.Context, user User) error
	Activate(ctx context.Context, userID string, token string) error
	GenerateToken(user User, tokentype string) (string, error)
	GenerateActivateToken(hashedpassword string) string
	GetProfile(ctx context.Context, id string) (Profile, error)
	ForgetPassword(ctx context.Context, email Email) error
	ResetPassword(ctx context.Context, userid, token, password, newPassword string) error
	GetUsers(ctx context.Context) ([]User, error)
	DeleteUser(ctx context.Context, id string) error
}
