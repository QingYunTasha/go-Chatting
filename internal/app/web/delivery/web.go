package delivery

import (
	webdomain "go-Chatting/domain/app/web"
	dbdomain "go-Chatting/domain/infra/database"
	"go-Chatting/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
	usecase webdomain.WebUsecase
}

func NewWebHandler(router *gin.Engine, webUsecase webdomain.WebUsecase, secretKey webdomain.SecretKey) {
	handler := &WebHandler{
		usecase: webUsecase,
	}
	// Add public routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.POST("/logout", handler.Logout)
	router.POST("/forgotpassword", handler.ForgotPassword)
	router.POST("/resetpassword", handler.ResetPassword)

	// Add authentication middleware to protected routes
	protectedRoutes := router.Group("/users")
	protectedRoutes.Use(AuthMiddleware(secretKey))
	protectedRoutes.GET("/:id", handler.ViewProfile)
	protectedRoutes.PATCH("/:id", handler.UpdateProfile)
	protectedRoutes.PATCH("/:id/password", handler.ChangePassword)
	protectedRoutes.POST("/:id/joingroup", handler.JoinGroup)
	protectedRoutes.POST("/:id/leavegroup", handler.LeaveGroup)
	protectedRoutes.POST("/:id/addfriend", handler.AddFriend)
	protectedRoutes.POST("/:id/removefriend", handler.RemoveFriend)

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
	if err := h.usecase.Logout(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user logged out successfully"})
}

func (h *WebHandler) ViewProfile(c *gin.Context) {
	userID := c.Param("id")

	// Convert string to uint32
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usecase.ViewProfile(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *WebHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("id")

	var user dbdomain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert string to uint32
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.UpdateProfile(uint32(id), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WebHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("id")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	err := h.usecase.ChangePassword(uint32(userID), oldPassword, newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *WebHandler) ForgotPassword(c *gin.Context) {
	email := c.PostForm("email")

	err := h.usecase.ForgotPassword(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (h *WebHandler) ResetPassword(c *gin.Context) {
	token := c.Query("token")
	password := c.PostForm("password")

	err := h.usecase.ResetPassword(token, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *WebHandler) JoinGroup(c *gin.Context) {
	userID := c.GetUint("id")
	groupID := c.PostForm("group_id")

	// Convert string to uint32
	groupIDUint32, err := utils.StringToUint32(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.usecase.JoinGroup(uint32(userID), groupIDUint32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined group successfully"})
}

func (h *WebHandler) LeaveGroup(c *gin.Context) {
	userID := c.GetUint("id")
	groupID := c.PostForm("group_id")

	// Convert string to uint32
	groupIDUint32, err := utils.StringToUint32(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.usecase.LeaveGroup(uint32(userID), groupIDUint32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left group successfully"})
}

func (h *WebHandler) AddFriend(c *gin.Context) {
	userID := c.GetUint("id")
	friendIDStr := c.PostForm("friend_id")

	// Convert string to uint32
	friendIDUint32, err := utils.StringToUint32(friendIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.usecase.AddFriend(uint32(userID), friendIDUint32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

func (h *WebHandler) RemoveFriend(c *gin.Context) {
	userID := c.GetUint("id")
	friendIDStr := c.PostForm("friend_id")

	// Convert string to uint32
	friendIDUint32, err := utils.StringToUint32(friendIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.usecase.RemoveFriend(uint32(userID), friendIDUint32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed successfully"})
}
