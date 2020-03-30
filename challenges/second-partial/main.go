package main

import (
	_ "image/jpeg"
	_ "image/png"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
)

type LoginStruct struct {
	User     string
	Password string
	Token    string
}

var db = make(map[string]LoginStruct)

func main() {
	router := gin.Default()

	router.GET("/login", LoginHandler)
	router.GET("/logout", LogoutHandler)
	router.GET("/status", StatusHandler)
	router.POST("/upload", UploadHandler)

	router.Run(":8080")
}

// Login User handler
func LoginHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	username := pair[0]
	password := pair[1]
	taken := false
	for i, _ := range db {
		if db[i].User == username {
			taken = true
		}
	}

	if taken || username == "" {
		c.JSON(http.StatusOK, ErrorMessageResponse("This username is taken"))
	} else {
		tokenString := auth[1]
		db[tokenString] = LoginStruct{
			User:     username,
			Password: password,
			Token:    tokenString,
		}
		c.JSON(http.StatusOK, SuccessLoginResponse(username, tokenString))
	}

}

// Log Out handler
func LogoutHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := db[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		username := db[token].User
		c.JSON(http.StatusOK, SuccessLogoutResponse(username))
		delete(db, token)
	}

}

// Status Handler Login User handler
func StatusHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := db[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		username := db[token].User
		c.JSON(http.StatusOK, SuccessStatusResponse(username))
	}
}

// Upload handler
func UploadHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := db[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		file, err := c.FormFile("data")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		f, err := os.Open(filename)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMessageResponse("There was an error with the image"))
		}
		image, _, err := image.DecodeConfig(f)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMessageResponse("There was an error opening image"))
		} else {
			c.JSON(http.StatusOK, SuccessUploadResponse(filename, image.Width, image.Height))
		}
	}

}

func SuccessLoginResponse(username string, token string) gin.H {
	return gin.H{
		"message": "Hi " + username + ", welcome to the DPIP System",
		"token":   token,
	}
}
func SuccessLogoutResponse(username string) gin.H {
	return gin.H{
		"message": "Bye " + username + ", your token has been revoked",
	}
}

// ErrorMessageResponse Request response object ready for errors.
func ErrorMessageResponse(message string) gin.H {
	return gin.H{
		"status": "error",
		"data": gin.H{
			"message": message,
		},
	}
}
func SuccessStatusResponse(username string) gin.H {
	return gin.H{
		"message": "Hi " + username + ", the DPIP System is Up and Running",
		"time":    time.Now().Format("2006-01-02T15:04:05+07:00"),
	}
}
func SuccessUploadResponse(image string, width int, height int) gin.H {
	return gin.H{
		"message": "Image: " + image + " uploaded succefully",
		"time":    time.Now().Format("2006-01-02T15:04:05+07:00"),
		"size":    strconv.Itoa(width) + "x" + strconv.Itoa(height),
	}
}

