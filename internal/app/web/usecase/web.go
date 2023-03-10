package usecase

import (
	"errors"
	"fmt"
	webdomain "go-Chatting/domain/app/web"
	dbdomain "go-Chatting/domain/infra/database"
	dbfactory "go-Chatting/internal/infra/database/factory"
	utils "go-Chatting/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type WebUsecase struct {
	DbRepo    *dbfactory.DatabaseRepository
	Mailer    webdomain.SMTPMailer
	SecretKey webdomain.SecretKey
}

func NewWebUsecase(dbRepo *dbfactory.DatabaseRepository, mailer webdomain.SMTPMailer, secretKey webdomain.SecretKey) webdomain.WebUsecase {
	return &WebUsecase{
		DbRepo:    dbRepo,
		Mailer:    mailer,
		SecretKey: secretKey,
	}
}

func (u *WebUsecase) Register(name, email, password string) error {
	// Validate input
	if name == "" || email == "" || password == "" {
		return errors.New("name, email, and password are required")
	}

	// Check if user already exists
	_, err := u.DbRepo.User.GetByEmail(email)
	if err == nil {
		return errors.New("email is already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Save user to repository
	user := &dbdomain.User{Name: name, Email: email, Password: string(hashedPassword)}
	return u.DbRepo.User.Create(user)
}

func (u *WebUsecase) Login(email, password string, w http.ResponseWriter) (uint32, error) {
	// Find user by email
	user, err := u.DbRepo.User.GetByEmail(email)
	if err != nil {
		return 0, errors.New("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("invalid email or password")
	}

	// Create session token
	token := jwt.New(jwt.SigningMethodES256)
	claims := webdomain.CustomJWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token.Claims = claims
	signedToken, err := token.SignedString([]byte(u.SecretKey.Get()))
	if err != nil {
		return 0, err
	}

	// Set cookie
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "token", Value: signedToken, Expires: expiration}
	http.SetCookie(w, &cookie)

	return user.ID, nil
}

func (u *WebUsecase) Logout(w http.ResponseWriter, r *http.Request) error {
	// Get cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return errors.New("not logged in")
	}

	// Delete cookie
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)

	return nil
}

func (u *WebUsecase) ViewProfile(userID uint32) (*dbdomain.User, error) {
	// Find user by ID
	user, err := u.DbRepo.User.Get(userID)
	if err != nil {
		return nil, err
	}

	// Mask password
	user.Password = ""

	return user, nil
}

func (u *WebUsecase) UpdateProfile(userID uint32, user *dbdomain.User) error {
	// Find user by ID
	oldUser, err := u.DbRepo.User.Get(userID)
	if err != nil {
		return err
	}

	// Update fields
	oldUser.Name = user.Name
	oldUser.Email = user.Email

	// Save changes
	err = u.DbRepo.User.Update(userID, oldUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) ChangePassword(userID uint32, oldPassword string, newPassword string) error {
	// Find user by ID
	user, err := u.DbRepo.User.Get(userID)
	if err != nil {
		return err
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	err = u.DbRepo.User.Update(userID, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) ForgotPassword(email string) error {
	// Find user by email
	user, err := u.DbRepo.User.GetByEmail(email)
	if err != nil {
		return err
	}

	// Generate password reset token
	token := utils.GenerateRandomToken(32)
	err = u.DbRepo.User.SavePasswordResetToken(user.ID, token)
	if err != nil {
		return err
	}

	// Send password reset email
	err = u.Mailer.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) ResetPassword(token string, password string) error {
	// Find user by password reset token
	user, err := u.DbRepo.User.FindByPasswordResetToken(token)
	if err != nil {
		return err
	}

	// Reset password
	err = u.DbRepo.User.ResetPassword(user.ID, password)
	if err != nil {
		return err
	}

	// Clear password reset token
	err = u.DbRepo.User.ClearPasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) CreateGroup(userID uint32, groupName string) error {
	group := &dbdomain.Group{
		Name:    groupName,
		OwnerID: userID,
	}

	return u.DbRepo.Group.Create(group)
}

func (u *WebUsecase) RemoveGroup(userID uint32, groupName string) error {
	group, err := u.DbRepo.Group.GetByName(groupName)
	if err != nil {
		return err
	}

	if group.OwnerID != userID {
		return fmt.Errorf("%s is not owned by %d", groupName, userID)
	}

	return u.DbRepo.Group.Delete(group.ID)
}

func (u *WebUsecase) JoinGroup(userID uint32, groupName string) error {
	// Check if user is already a member of the group
	isMember, err := u.DbRepo.User.IsGroupMember(userID, groupName)
	if err != nil {
		return err
	}

	if isMember {
		// User is already a member of the group
		return errors.New("user is already a member of the group")
	}

	// Add user to group
	err = u.DbRepo.User.AddGroupMember(userID, groupName)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) LeaveGroup(userID uint32, groupName string) error {
	// Check if user is a member of the group
	isMember, err := u.DbRepo.User.IsGroupMember(userID, groupName)
	if err != nil {
		return err
	}

	if !isMember {
		// User is not a member of the group
		return errors.New("user is not a member of the group")
	}

	// Remove user from group
	err = u.DbRepo.User.RemoveGroupMember(userID, groupName)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) AddFriend(userID uint32, friendEmail string) error {
	// Check if the user is already friends with the friend
	areFriends, err := u.DbRepo.User.AreFriends(userID, friendEmail)
	if err != nil {
		return err
	}

	if areFriends {
		// The user is already friends with the friend
		return errors.New("the user is already friends with the friend")
	}

	// Add the friend to the user's friend list
	err = u.DbRepo.User.AddFriend(userID, friendEmail)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) RemoveFriend(userID uint32, friendEmail string) error {
	// Check if the user is friends with the friend
	areFriends, err := u.DbRepo.User.AreFriends(userID, friendEmail)
	if err != nil {
		return err
	}

	if !areFriends {
		// The user is not friends with the friend
		return errors.New("the user is not friends with the friend")
	}

	// Remove the friend from the user's friend list
	err = u.DbRepo.User.RemoveFriend(userID, friendEmail)
	if err != nil {
		return err
	}

	return nil
}
