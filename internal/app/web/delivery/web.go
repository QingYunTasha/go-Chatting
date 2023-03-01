package delivery

import (
	"errors"
	webdomain "go-Chatting/domain/app/web"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
	usecase webdomain.WebUsecase
}

func NewWebHandler(s *gin.Engine, webUsecase webdomain.WebUsecase) {
	h := &WebHandler{
		usecase: webUsecase,
	}

	s.POST("/register", h.Register)
	s.GET("/reset-password", h.ResetPassword)
}

func (h *WebHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.Register(req.Name, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *WebHandler) Login(c *gin.Context) {
	// Parse email and password from request body
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attempt to log in the user
	err := h.usecase.Login(req.Email, req.Password, c.Writer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WebHandler) Logout(c *gin.Context) {
	if err := d.usecase.Logout(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to logout user").Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user logged out successfully"})
}

func (h *WebHandler) ViewProfile(c *gin.Context) {
	userID := c.Param("userID")

	// Convert string to uint32
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "invalid user ID").Error()})
		return
	}

	user, err := d.usecase.ViewProfile(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "failed to get user profile").Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *WebHandler) UpdateProfile(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	var user dbdomain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.UpdateProfile(userID, &user); err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	h.respondSuccess(c)
}

func (h *WebHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint32("userID")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	err := h.UseCase.ChangePassword(userID, oldPassword, newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *WebHandler) ForgotPassword(c *gin.Context) {
	email := c.PostForm("email")

	err := h.UseCase.ForgotPassword(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (h *WebHandler) ResetPassword(c *gin.Context) {
	token := c.Query("token")
	password := c.PostForm("password")

	err := h.UseCase.ResetPassword(token, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *WebHandler) JoinGroup(c *gin.Context) {
	userID := c.GetUint32("userID")
	groupID := c.PostForm("group_id")

	err := h.UseCase.JoinGroup(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined group successfully"})
}

func (h *WebHandler) LeaveGroup(c *gin.Context) {
	userID := c.GetUint32("userID")
	groupID := c.PostForm("group_id")

	err := h.UseCase.LeaveGroup(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left group successfully"})
}

func (h *WebHandler) AddFriend(c *gin.Context) {
	userID := c.GetUint32("userID")
	friendID := c.PostForm("friend_id")

	err := h.UseCase.AddFriend(userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

func (h *WebHandler) RemoveFriend(c *gin.Context) {
	userID := c.GetUint32("userID")
	friendID := c.PostForm("friend_id")

	err := h.UseCase.RemoveFriend(userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed successfully"})
}
