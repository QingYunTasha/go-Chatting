package domain

import (
	dbdomain "go-Chatting/domain/infra/database"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type WebUsecase interface {
	Register(name, email, password string) error
	Login(email, password string, w http.ResponseWriter) error
	Logout(w http.ResponseWriter, r *http.Request) error
	ViewProfile(userID uint32) (*dbdomain.User, error)
	UpdateProfile(userID uint32, user *dbdomain.User) error
	ChangePassword(userID uint32, oldPassword string, newPassword string) error
	ForgotPassword(email string) error
	ResetPassword(token string, password string) error
	JoinGroup(userID uint32, groupID uint32) error
	LeaveGroup(userID uint32, groupID uint32) error
	AddFriend(userID uint32, friendID uint32) error
	RemoveFriend(userID uint32, friendID uint32) error
}

type CustomJWTClaims struct {
	UserID uint32 `json:"user_id"`
	jwt.RegisteredClaims
}

type SMTPMailer interface {
	SendPasswordResetEmail(to string, token string) error
}

type SecretKey interface {
	Get() string
}
