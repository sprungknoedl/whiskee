package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func LogoutR(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func CallbackR(c *gin.Context) {
	// dependency injection
	var (
		userStore = c.MustGet(UserStoreKey).(UserStore)
	)

	domain := "whiskee.eu.auth0.com"
	redirect := url.URL{
		Scheme: "http",
		Host:   c.Request.Host,
		Path:   "/auth/callback",
	}

	// Instantiating the OAuth2 package to exchange the Code for a Token
	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CID"),
		ClientSecret: os.Getenv("AUTH0_SECRET"),
		RedirectURL:  redirect.String(),
		Scopes:       []string{"openid", "name", "email", "nickname"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	// Getting the Code that we got from Auth0
	code := c.Query("code")

	// Exchanging the code for a token
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Getting now the User information
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Unmarshalling the JSON of the Profile
	var profile map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	resp.Body.Close()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user := User{
		Auth0:   profile["user_id"].(string),
		Name:    profile["name"].(string),
		Email:   profile["email"].(string),
		Picture: profile["picture"].(string),
	}

	// store user in database
	user, err = userStore.SaveUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// store user in session
	session := sessions.Default(c)
	session.Set("user", user)
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}
