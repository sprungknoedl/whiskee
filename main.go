package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/dustin/go-humanize"
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
	router.SetHTMLTemplate(template.Must(template.New("views").
		Funcs(template.FuncMap{
		"humanizeTime": humanize.Time,
	}).
		ParseGlob("templates/*")))

	router.GET("/", IndexR)
	router.GET("/auth/google", GoogleAuthR)
	router.GET("/auth/google/callback", GoogleCallbackR)

	secured := router.Group("/", Secured)
	secured.GET("/logout", LogoutR)

	secured.GET("/feed", NewsFeedR)
	secured.GET("/users", UsersR)
	secured.GET("/friends", FriendsR)
	secured.GET("/u/:id", ProfileR)
	secured.GET("/u/:id/friend", AddFriendR)

	secured.POST("/p", AddPostR)
	secured.POST("/p/:id/comment", CommentR)
	secured.GET("/p/:id/delete", DeletePostR)

	secured.POST("/w", AddWhiskeyR)

	router.Run(":" + port)
}

func IndexR(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func NewsFeedR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)

	whiskeys, err := db.Whiskeys.All()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	posts, err := db.Posts.GetNewsFeed(principal.ID, 50)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"principal": principal,
		"whiskeys":  whiskeys,
		"posts":     posts,
	})
}

func ProfileR(c *gin.Context) {
	// principal := c.MustGet("user").(*model.User)
	// id := c.Param("id")

	// user, err := db.Users.Get(id)
	// if err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }
	c.Redirect(http.StatusSeeOther, "/feed")
}

type WhiskeyForm struct {
	Distillery string  `form:"distillery"`
	Name       string  `form:"name"`
	Type       string  `form:"type"`
	Age        int     `form:"age"`
	ABV        float64 `form:"abv"`
	Size       float64 `form:"size"`
}

func AddWhiskeyR(c *gin.Context) {
	var form WhiskeyForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	whiskey := &model.Whiskey{
		Distillery: form.Distillery,
		Name:       form.Name,
		Type:       form.Type,
		Age:        form.Age,
		ABV:        form.ABV,
		Size:       form.Size,
	}

	if _, err := db.Whiskeys.Create(whiskey); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/feed")
}

type PostForm struct {
	Whiskey int    `form:"whiskey"`
	Body    string `form:"body"`
}

func AddPostR(c *gin.Context) {
	principal := c.MustGet("user").(*model.User)

	var form PostForm
	if err := c.Bind(&form); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post := &model.Post{
		User:    principal,
		Date:    time.Now(),
		Whiskey: &model.Whiskey{ID: form.Whiskey},
		Body:    form.Body,
	}

	if _, err := db.Posts.Create(post); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/u/"+principal.ID)
}

func CommentR(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/feed")
}

func DeletePostR(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := db.Posts.Delete(id); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/feed")
}
