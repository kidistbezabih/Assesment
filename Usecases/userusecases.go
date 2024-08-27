package auth

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	domain "github.com/kidistbezabih/loan-tracker-api/Domain"
	infrastructure "github.com/kidistbezabih/loan-tracker-api/Infrastructure"
	"github.com/kidistbezabih/loan-tracker-api/errs"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserUsecase struct {
	repository   domain.AuthRepository
	emailService infrastructure.EmailService
}

func NewAuthUserUsecase(repository domain.AuthRepository, emailService infrastructure.EmailService) domain.AuthServices {
	return &AuthUserUsecase{
		repository:   repository,
		emailService: emailService,
	}
}

func (au *AuthUserUsecase) Login(ctx context.Context, info domain.LoginForm) (string, string, error) {
	user, err := au.repository.GetUserByUsername(ctx, info.Username)
	if err != nil {
		return "", "", errs.ErrNoUserWithUsername
	}
	if !user.IsActive {
		return "", "", errs.ErrAccountNotActivated
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password))
	if err != nil {
		return "", "", errs.ErrIncorrectPassword
	}

	refToken, err := au.GenerateToken(user, "refresh")
	if err != nil {
		return "", "", err
	}
	accessToken, err := au.GenerateToken(user, "access")
	if err != nil {
		return "", "", err
	}
	err = au.repository.RegisterRefreshToken(ctx, user.ID, refToken)
	if err != nil {
		return "", "", err
	}

	return refToken, accessToken, nil
}

func (au *AuthUserUsecase) RegisterUser(ctx context.Context, user domain.User) error {
	// var newUser User
	_, err := au.repository.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return errs.ErrUserExistWithThisEmail
	}
	_, err = au.repository.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return errs.ErrUserExistWithThisUsername
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedpassword)
	user.IsActive = false
	user.Email = strings.ToLower(user.Email)
	// user.IsAdmin = false	activationLink := fmt.Sprintf("http://localhost/activate/%s/%s", user.ID, tokenString)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	// if the user is first user is make it admin and super admin
	count, err := au.repository.GetCollectionCount(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		user.IsAdmin = true
	}
	id, err := au.repository.CreateUser(ctx, user)
	if err != nil {
		return errs.ErrCantCreateUser
	}
	user.ID = id

	from := os.Getenv("FROM")
	tokenString := au.GenerateActivateToken(user.Password, user.UpdatedAt)

	activationLink := fmt.Sprintf("http://localhost/activate/%s/%s", user.ID, tokenString)
	au.emailService.SendEmail(from, user.Email, fmt.Sprintf("click the link to activate you account %s", activationLink), "Account Activation")

	return nil
}

func (au *AuthUserUsecase) GetProfile(ctx context.Context, id string) (domain.Profile, error) {
	var profile domain.Profile
	user, err := au.repository.GetUserByID(ctx, id)
	if err != nil {
		return domain.Profile{}, err
	}
	profile.Email = user.Email
	profile.Username = user.Username
	profile.Name = user.Name
	return profile, err
}

func (au *AuthUserUsecase) Activate(ctx context.Context, userID string, token string) error {
	user, err := au.repository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	expectedToken := au.GenerateActivateToken(user.Password, user.UpdatedAt)

	if expectedToken != token {
		return err
	}

	user.IsActive = true
	user.UpdatedAt = time.Now()

	_, err = au.repository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthUserUsecase) GenerateActivateToken(hashedpassword string, updatedat time.Time) string {
	token := hashedpassword + updatedat.String()
	hasher := sha1.New()
	hasher.Write([]byte(token))

	token = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return token
}

func (au *AuthUserUsecase) GenerateToken(user domain.User, tokenType string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	claims := jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
		"isadmin":  user.IsAdmin,
		"isactive": user.IsActive,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"type":     tokenType,
	}
	if tokenType == "refresh" {
		claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (au *AuthUserUsecase) ForgetPassword(ctx context.Context, email domain.Email) error {
	user, err := au.repository.GetUserByEmail(ctx, email.User_email)
	if err != nil {
		return errs.ErrNoUserWithEmail
	}
	currenttime := time.Now().String()
	token := au.GenerateTokenForReset(ctx, user.UpdatedAt.String(), user.Password, currenttime)

	URLSafe := base64.URLEncoding.EncodeToString([]byte(currenttime))
	// send the token to that email
	from := os.Getenv("FROM")
	link := fmt.Sprintf("http://localhost:8000/v1/auth/reset/%s/%s/%s", user.ID, URLSafe, token)
	au.emailService.SendEmail(from, email.User_email, fmt.Sprintf("click the link to activate your password %s ", link), "Reset password")
	return nil
}

func (au *AuthUserUsecase) GenerateTokenForReset(ctx context.Context, updatedat, hashedpassword, currenttime string) string {
	data := hashedpassword + updatedat + currenttime
	hash := sha256.New()
	hash.Write([]byte(data))
	token := hex.EncodeToString(hash.Sum(nil))

	return token
}

func (au *AuthUserUsecase) ResetPassword(ctx context.Context, userid, tokenTime, token, password, newPassword string) error {
	user, _ := au.repository.GetUserByID(ctx, userid)

	expectedToken := au.GenerateTokenForReset(ctx, user.UpdatedAt.String(), user.Password, tokenTime)
	if expectedToken != token {
		return errors.New("some error")
	}
	return nil
}
