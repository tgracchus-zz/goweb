package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string
}

type UserAction struct {
	User   User
	Action string
}

func main() {
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		hello := User{c.Param("name")}
		c.JSON(http.StatusOK, hello)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		userAction := UserAction{User{name}, action}
		c.JSON(http.StatusOK, userAction)
	})

	router.Run(":8080")
}
