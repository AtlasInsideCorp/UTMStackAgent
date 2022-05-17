package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type changePasswordJSON struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func login(c *gin.Context) {
	var json loginJSON
	session := sessions.Default(c)

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := getUser(hostname)
	if (user == User{}) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	match := CheckPasswordHash(json.Password, user.password)
	if match {
		//session.Set("hostname", user.ip)
		session.Set("id", user.id)
		session.Save()

		insertLog(LogObj{Id: 0, LogType: "Login", Desc: hostname + " loged in", Date: time.Now()})
		c.JSON(http.StatusOK, gin.H{"success": true})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	//session.Set("id", "")
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()

	insertLog(LogObj{Id: 0, LogType: "LogOut", Desc: hostname + " signed out", Date: time.Now()})

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func changePassword(c *gin.Context) {
	var json changePasswordJSON

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := getUser(hostname)

	match := CheckPasswordHash(json.Password, user.password)
	if match {
		hash, _ := HashPassword(json.NewPassword)
		userChangePassword(hostname, hash)

		insertLog(LogObj{Id: 0, LogType: "ChangePassword", Desc: hostname + " changed password", Date: time.Now()})
		c.JSON(http.StatusOK, gin.H{"success": true})
	} else {
		insertLog(LogObj{Id: 0, LogType: "ChangePassword", Desc: hostname + " changed password error", Date: time.Now()})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please review your current password"})
	}
}

func createPassword(c *gin.Context) {
	var json loginJSON

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := getUser(hostname)
	if (user != User{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already have a password"})
		return
	}

	hash, _ := HashPassword(json.Password)
	insertUser(hostname, hash)
	insertLog(LogObj{Id: 0, LogType: "SetPassword", Desc: hostname + " created first password", Date: time.Now()})
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func saveSettings(c *gin.Context) {
	var json Settings

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings := getSettings()
	if (settings == Settings{}) {
		json.Date = time.Now()
		insertSettings(json)

		insertLog(LogObj{Id: 0, LogType: "CreateSettings", Desc: hostname + " created first settings", Date: time.Now()})
	} else {
		json.Date = time.Now()
		updateSettings(json)
		insertLog(LogObj{Id: 0, LogType: "UpdateSettings", Desc: hostname + " updated settings", Date: time.Now()})
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func editModule(c *gin.Context) {
	var json Module

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modules := getModules()
	if len(modules) > 0 {
		m := Module{}
		for i := range modules {
			if modules[i].Id == json.Id {
				m = modules[i]
				break
			}
		}

		json.Date = time.Now()
		json.Name = m.Name
		json.Image = m.Image

		updateModule(json)

		insertLog(LogObj{Id: 0, LogType: "UpdateModule", Desc: hostname + " updated module " + json.Name, Date: time.Now()})
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
