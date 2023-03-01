package usecase

import (
	"errors"
	dbdomain "go-Chatting/domain/infra/database"
	dbfactory "go-Chatting/internal/infra/database/factory"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type WebUsecase struct {
	DbRepo *dbfactory.DatabaseRepository
}

func NewWebUsecase(dbRepo *dbfactory.DatabaseRepository) *WebUsecase {
	return &WebUsecase{DbRepo: dbRepo}
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

func (u *WebUsecase) Login(email, password string, w http.ResponseWriter) error {
	// Find user by email
	user, err := u.DbRepo.User.FindByEmail(email)
	if err != nil {
		return errors.New("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("invalid email or password")
	}

	// Create session token
	token := uuid.New().String()

	// Set cookie
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
	http.SetCookie(w, &cookie)

	return nil
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
	token := generateRandomToken(32)
	err = u.savePasswordResetToken(user.ID, token)
	if err != nil {
		return err
	}

	// Send password reset email
	err = u.mailer.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomToken(length int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (u *WebUsecase) ResetPassword(token string, password string) error {
	// Find user by password reset token
	user, err := u.findByPasswordResetToken(token)
	if err != nil {
		return err
	}

	// Reset password
	err = u.resetPassword(user.ID, password)
	if err != nil {
		return err
	}

	// Clear password reset token
	err = u.clearPasswordResetToken(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) JoinGroup(userID uint32, groupID uint32) error {
	// Check if user is already a member of the group
	isMember, err := u.isGroupMember(userID, groupID)
	if err != nil {
		return err
	}

	if isMember {
		// User is already a member of the group
		return errors.New("user is already a member of the group")
	}

	// Add user to group
	err = u.addGroupMember(userID, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) LeaveGroup(userID uint32, groupID uint32) error {
	// Check if user is a member of the group
	isMember, err := u.isGroupMember(userID, groupID)
	if err != nil {
		return err
	}

	if !isMember {
		// User is not a member of the group
		return errors.New("user is not a member of the group")
	}

	// Remove user from group
	err = u.removeGroupMember(userID, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) AddFriend(userID uint32, friendID uint32) error {
	// Check if the user is already friends with the friend
	areFriends, err := u.areFriends(userID, friendID)
	if err != nil {
		return err
	}

	if areFriends {
		// The user is already friends with the friend
		return errors.New("the user is already friends with the friend")
	}

	// Add the friend to the user's friend list
	err = u.addFriend(userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (u *WebUsecase) RemoveFriend(userID uint32, friendID uint32) error {
	// Check if the user is friends with the friend
	areFriends, err := u.areFriends(userID, friendID)
	if err != nil {
		return err
	}

	if !areFriends {
		// The user is not friends with the friend
		return errors.New("the user is not friends with the friend")
	}

	// Remove the friend from the user's friend list
	err = u.removeFriend(userID, friendID)
	if err != nil {
		return err
	}

	return nil
}