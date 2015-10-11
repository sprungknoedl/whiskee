package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/contrib/sessions"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/kalaspuffar/base64url"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/golang.org/x/oauth2"
)

func Secured(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.Set("user", user)
	c.Next()
}

func LogoutR(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}

func GoogleAuthR(c *gin.Context) {
	state := "state"
	url := sso.AuthCodeURL(state)

	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	c.Redirect(http.StatusSeeOther, url)
}

func GoogleCallbackR(c *gin.Context) {
	session := sessions.Default(c)
	code := c.Query("code")
	state := c.Query("state")
	state2 := session.Get("state")

	if state != state2 {
		err := errors.New("invalid state")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tok, err := sso.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jwt, err := DecodeJWT(tok.Extra("id_token").(string))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := jwt["sub"].(string)
	email := jwt["email"].(string)

	user, err := db.Users.GetOrCreate(id, email)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set("user", user)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/feed")
}

type JWT map[string]interface{}

func DecodeJWT(token string) (JWT, error) {
	parts := strings.SplitN(token, ".", 3)
	jwt := make(map[string]interface{})

	payload, err := base64url.Decode(parts[1])
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(payload, &jwt); err != nil {
		return nil, err
	}

	return jwt, nil
}
