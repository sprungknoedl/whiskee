package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/contrib/commonlog"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/contrib/cors"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/contrib/sessions"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/model"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/golang.org/x/oauth2/google"
)

var (
	sso *oauth2.Config
	db  *model.Model
)

func main() {
	var err error
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
	store := sessions.NewCookieStore([]byte(key))

	dburl := os.Getenv("DATABASE_URL")
	if db, err = model.Connect(dburl); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	router := gin.New()
	router.Use(commonlog.New())
	router.Use(cors.Default())
	router.Use(sessions.Sessions("whiskee", store))
	router.Use(gin.Recovery())
	router.Use(gin.ErrorLogger())

	router.Static("/assets", "assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", IndexR)
	router.GET("/auth/google", GoogleAuthR)
	router.GET("/auth/google/callback", GoogleCallbackR)

	secured := router.Group("/", Secured)
	secured.GET("/u/:id", ProfileR)
	secured.GET("/logout", LogoutR)

	router.Run(":" + port)
}

func IndexR(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func ProfileR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)
	id := c.Param("id")

	user, err := db.Users.Get(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	whiskeys, err := db.Whiskeys.All()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	posts, err := db.Posts.FindByUser(id, 50)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"principal": principal,
		"user":      user,
		"whiskeys":  whiskeys,
		"posts":     posts,
		"home":      user.ID == principal.ID,
	})
}
