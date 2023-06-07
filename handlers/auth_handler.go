package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todo-api/models"
	"todo-api/utils"
)

func SignIn(c *gin.Context) {
	var request struct {
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := models.GetUserByEmail(request.EmailAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.ComparePasswords(user.Password, request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate and return an access token
	token, err := utils.GenerateAccessToken(user.UserID, user.EmailAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SignUp(c *gin.Context) {
	var request struct {
		Name         string `json:"name"`
		PhoneNumber  int    `json:"phone_number"`
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the Email is already taken
	existingUser, err := models.GetUserByEmail(request.EmailAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Create a new user
	user := &models.User{
		Name:         request.Name,
		PhoneNumber:  request.PhoneNumber,
		EmailAddress: request.EmailAddress,
		Password:     utils.HashPassword(request.Password),
	}

	// Save the user to the database
	if _, err := models.CreateUser(user.Name, user.PhoneNumber, user.EmailAddress, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
