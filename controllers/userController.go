package controllers

import (
	"example/REST_API_JWT/db"
	"example/REST_API_JWT/models"
	// "log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read the body",
			"success": false,
		})
		return
	}

	// hash the password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Hash the password",
			"success": false,
		})
		return
	}

	// create the user now
	newUser := models.User{
		Email:    body.Email,
		Passowrd: string(hashedPassword),
	}
	result := db.DB.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error Creating new user",
			"success": false,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created new user",
		"success": true,
	})
}

func SingIn(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read the body",
			"success": false,
		})
		return
	}

	// find the user if that exists on the server
	// var existUser models.User
	// db.DB.First(&existUser, "email = ?", body.Email)

	var existUser = models.User{Email: body.Email}
	db.DB.First(&existUser)

	if existUser.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
			"success": false,
		})
	}
	// check the hashed password mathches the password
	err := bcrypt.CompareHashAndPassword([]byte(existUser.Passowrd), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
			"success": false,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existUser.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// set it as cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"success": true,
		"token":   tokenString,
	})

}

func GetUserDetails(c *gin.Context) {
	id, _ := c.Get("id")
	var existUser models.User
	query := "SELECT ID, Email FROM users WHERE ID = ?"
	db.DB.Raw(query, id).Scan(&existUser)
	// db.DB.First(&existUser, id).Select("-Email")
	if existUser.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User feteched successfully",
		"success": true,
		"data":   existUser,
	})
}
