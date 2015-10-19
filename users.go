package main

import (
	"log"
	"net/http"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/model"
)

func SearchR(c *gin.Context) {
	query := c.Query("q")
	log.Printf("searching for %q", query)

	users, err := db.Users.SearchByEMail(query)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"items":  users,
	})
}

func UsersR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)

	users, err := db.Users.All()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "users.html", gin.H{
		"principal": principal,
		"users":     users,
	})
}

func FriendsR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)

	users, err := db.Users.Friends(principal.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "users.html", gin.H{
		"principal": principal,
		"users":     users,
	})
}

func AddFriendR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)
	friend := c.Param("id")

	if err := db.Users.AddFriend(principal.ID, friend); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/feed")
}

func DelFriendR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)
	friend := c.Param("id")

	if err := db.Users.DelFriend(principal.ID, friend); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/feed")
}
