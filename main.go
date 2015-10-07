package main

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gorilla/sessions"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/kalaspuffar/base64url"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/golang.org/x/oauth2/google"
)

var (
	sso   *oauth2.Config
	store sessions.Store
)

func init() {
	gob.Register(&User{})
}

func main() {
	port := os.Getenv("PORT")
	base := os.Getenv("BASE_URL")

	cid := os.Getenv("GOOGLE_CID")
	secret := os.Getenv("GOOGLE_SECRET")
	sso = &oauth2.Config{
		ClientID:     cid,
		ClientSecret: secret,
		RedirectURL:  base + "/auth/google/callback",
		Scopes:       []string{"openid", "email"},
		Endpoint:     google.Endpoint,
	}

	key := os.Getenv("SESSION_SECRET")
	store = sessions.NewCookieStore([]byte(key))

	router := gin.Default()
	router.Static("/assets", "assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", IndexR)
	router.GET("/auth/google", GoogleAuthR)
	router.GET("/auth/google/callback", GoogleCallbackR)

	secured := router.Group("/", Secured)
	secured.GET("/dashboard", DashboardR)
	secured.GET("/logout", LogoutR)

	router.Run(":" + port)
}

type User struct {
	ID       string
	EMail    string
	Gravatar string
}

func NewUser(id, email string) *User {
	data := []byte(strings.TrimSpace(email))
	gravatar := fmt.Sprintf("%x", md5.Sum(data))

	return &User{
		ID:       id,
		EMail:    email,
		Gravatar: gravatar,
	}
}

func Secured(c *gin.Context) {
	session, _ := store.Get(c.Request, "whiskee")
	user, ok := session.Values["user"]

	fmt.Printf("session = %v\n", session.Values)

	if !ok {
		fmt.Printf("not authenticated!\n")
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.Set("user", user)
	c.Next()
}

func IndexR(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func DashboardR(c *gin.Context) {
	user := c.MustGet("user").(*User)
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user": user,
	})
}

func LogoutR(c *gin.Context) {
	session, _ := store.Get(c.Request, "whiskee")
	session.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	c.Redirect(http.StatusSeeOther, "/")
}

func GoogleAuthR(c *gin.Context) {
	state := "state"
	url := sso.AuthCodeURL(state)

	session, _ := store.Get(c.Request, "whiskee")
	session.Values["state"] = state

	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusSeeOther, url)
}

func GoogleCallbackR(c *gin.Context) {
	session, _ := store.Get(c.Request, "whiskee")
	code := c.Query("code")
	state := c.Query("state")
	state2 := session.Values["state"]

	if state != state2 {
		c.String(http.StatusBadRequest, "invalid state %q != %q", state, state2)
		return
	}

	tok, err := sso.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	jwt, err := DecodeJWT(tok.Extra("id_token").(string))
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("successfully authenticated user\n")
	session.Values["user"] = NewUser(jwt["sub"].(string), jwt["email"].(string))
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusSeeOther, "/dashboard")
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
