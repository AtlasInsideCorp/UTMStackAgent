package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginJSON struct {
	Password string `json:"password"`
}

func isAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	sessionID := session.Get("id")
	if sessionID == nil {
		return false
	}

	return true
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c) {
			/*c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})*/
			log.Println("redirect")
			c.Redirect(http.StatusMovedPermanently, "/sign-in")
			c.Abort()

		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
