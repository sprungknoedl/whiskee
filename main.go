package main

import (
	_ "crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/transloadit/go-sdk"
)

const LIMIT = 8
const TransloaditKey = "github.com/sprungknoedl/whiskee/transloadit-client"

func main() {
	log.Printf("setting up gin ...")
	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")

	log.Printf("setting up session store ...")
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("whiskee", store))

	log.Printf("setting up database ...")
	r.Use(PostgresDatabase())

	log.Printf("setting up transloadit ...")
	r.Use(Transloadit())

	log.Printf("setting up routes ...")
	r.GET("/", HomeR)
	r.GET("/auth/callback", CallbackR)
	r.GET("/auth/logout", LogoutR)

	r.GET("/home", HomeR)
	r.GET("/whisky", WhiskyListR)
	r.GET("/whisky/:id", WhiskyR)

	r.GET("/add/whisky", AddWhiskyFormR)
	r.POST("/add/whisky", AddWhiskyR)
	r.GET("/edit/whisky/:id", EditWhiskyFormR)
	r.POST("/edit/whisky/:id", EditWhiskyR)

	r.POST("/add/review", AddReviewR)
	r.GET("/edit/review/:id", EditReviewFormR)
	r.POST("/edit/review/:id", EditReviewR)

	port := os.Getenv("PORT")
	log.Printf("listening on :%s", port)
	r.Run(":" + port)
}

func Transloadit() gin.HandlerFunc {
	options := transloadit.DefaultConfig
	options.AuthKey = os.Getenv("TRANSLOADIT_AUTH_KEY")
	options.AuthSecret = os.Getenv("TRANSLOADIT_SECRET_KEY")
	client, err := transloadit.NewClient(options)
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set(TransloaditKey, client)
		c.Next()
	}
}

func GetUser(c *gin.Context) User {
	session := sessions.Default(c)
	if obj := session.Get("user"); obj != nil {
		return obj.(User)
	}

	return User{}
}

func HomeR(c *gin.Context) {
	// dependency injection
	var (
		whiskyStore = c.MustGet(WhiskyStoreKey).(WhiskyStore)
		reviewStore = c.MustGet(ReviewStoreKey).(ReviewStore)
		user        = GetUser(c)
	)

	// temporary variables
	var (
		activity []Review
		trending []Whisky
		err      error
	)

	// get activity stream for user
	activity, err = reviewStore.GetActivity(user.ID, LIMIT)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// get trending whiskies
	trending, err = whiskyStore.GetTrending(LIMIT)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"User":     user,
		"Activity": activity,
		"Trending": trending,
	})
}

func WhiskyListR(c *gin.Context) {
	// dependency injection
	var (
		whiskyStore = c.MustGet(WhiskyStoreKey).(WhiskyStore)
		user        = GetUser(c)
	)

	// get whiskies
	list, err := whiskyStore.GetAllWhiskies()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// transform list into alphabet map
	whiskies := map[string][]Whisky{}
	for _, whisky := range list {
		group := strings.ToUpper(whisky.Distillery[0:1])
		whiskies[group] = append(whiskies[group], whisky)
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"User":     user,
		"Whiskies": whiskies,
	})
}

func WhiskyR(c *gin.Context) {
	// dependency injection
	var (
		whiskyStore = c.MustGet(WhiskyStoreKey).(WhiskyStore)
		reviewStore = c.MustGet(ReviewStoreKey).(ReviewStore)
		user        = GetUser(c)
	)

	// parse id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// fetch whiksy
	whisky, err := whiskyStore.GetWhisky(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// fetch reviews
	reviews, err := reviewStore.GetAllReviews(id, 30)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "whisky.html", gin.H{
		"User":    user,
		"Whisky":  whisky,
		"Reviews": reviews,
	})
}

func AddWhiskyFormR(c *gin.Context) {
	// dependency injection
	var (
		user = GetUser(c)
	)

	// only authenticated users can add a whisky
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.HTML(http.StatusOK, "whisky-form.html", gin.H{
		"Action": "Add",
		"User":   user,
		"Form":   Whisky{},
		"Errors": map[string]string{},
	})
}

func AddWhiskyR(c *gin.Context) {
	// dependency injection
	var (
		whiskyStore = c.MustGet(WhiskyStoreKey).(WhiskyStore)
		user        = GetUser(c)
	)

	// only authenticated users can perform this action
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	form := Whisky{}
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok, errors := ValidateWhisky(form)
	if !ok {
		c.HTML(http.StatusConflict, "whisky-form.html", gin.H{
			"Action": "Add",
			"User":   user,
			"Form":   form,
			"Errors": errors,
		})
		return
	}

	// store image
	var err error
	form, err = StoreWhiskyImage(c, form)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save whisky
	whisky, err := whiskyStore.SaveWhisky(form)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	url := fmt.Sprintf("/whisky/%d", whisky.ID)
	c.Redirect(http.StatusSeeOther, url)
}

func EditWhiskyFormR(c *gin.Context) {
	// dependency injection
	var (
		whiskyStore = c.MustGet(WhiskyStoreKey).(WhiskyStore)
		user        = GetUser(c)
	)

	// only authenticated users can perform this action
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// fetch whisky to edit
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	whisky, err := whiskyStore.GetWhisky(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "whisky-form.html", gin.H{
		"Action": "Edit",
		"User":   user,
		"Form":   whisky,
		"Errors": map[string]string{},
	})
}

func EditWhiskyR(c *gin.Context) {
	// dependency injection
	var (
		whiskyStore = c.MustGet(WhiskyStoreKey).(WhiskyStore)
		user        = GetUser(c)
	)

	// only authenticated users can perform this action
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// bind and validate form
	form := Whisky{}
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok, errors := ValidateWhisky(form)
	if !ok {
		c.HTML(http.StatusConflict, "whisky-form.html", gin.H{
			"Action": "Edit",
			"User":   user,
			"Form":   form,
			"Errors": errors,
		})
		return
	}

	// store image
	var err error
	form, err = StoreWhiskyImage(c, form)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// update whisky
	whisky, err := whiskyStore.SaveWhisky(form)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/whisky/%d", whisky.ID)
	c.Redirect(http.StatusSeeOther, url)
}

func StoreWhiskyImage(c *gin.Context, whisky Whisky) (Whisky, error) {
	var client = c.MustGet(TransloaditKey).(*transloadit.Client)

	file, header, err := c.Request.FormFile("file")
	if err == nil {
		assembly := client.CreateAssembly()
		assembly.TemplateId = os.Getenv("TRANSLOADIT_TEMPLATE_ID")
		assembly.Blocking = true
		assembly.AddReader("image", header.Filename, file)

		info, err := assembly.Upload()
		if err != nil {
			return whisky, err
		}

		whisky.Picture = info.Results[":original"][0].Url
		whisky.Thumbnail = info.Results["resize"][0].Url
	}

	return whisky, nil
}

func AddReviewR(c *gin.Context) {
	// dependency injection
	var (
		reviewStore = c.MustGet(ReviewStoreKey).(ReviewStore)
		user        = GetUser(c)
	)

	// only authenticated users can perform this action
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// bind and validate review
	form := Review{}
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	form.UserID = user.ID
	ok, errors := ValidateReview(form)
	if !ok {
		c.HTML(http.StatusConflict, "review-form.html", gin.H{
			"Action": "Add",
			"User":   user,
			"Form":   form,
			"Errors": errors,
		})
		return
	}

	// saving review
	review, err := reviewStore.SaveReview(form)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	url := fmt.Sprintf("/whisky/%d", review.Whisky.ID)
	c.Redirect(http.StatusSeeOther, url)
}

func EditReviewFormR(c *gin.Context) {
	// dependency injection
	var (
		reviewStore = c.MustGet(ReviewStoreKey).(ReviewStore)
		user        = GetUser(c)
	)

	// only authenticated users can perform this action
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// fetch review
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	review, err := reviewStore.GetReview(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "review-form.html", gin.H{
		"Action": "Edit",
		"User":   user,
		"Form":   review,
		"Errors": map[string]string{},
	})
}

func EditReviewR(c *gin.Context) {
	// dependency injection
	var (
		reviewStore = c.MustGet(ReviewStoreKey).(ReviewStore)
		user        = GetUser(c)
	)

	// only authenticated users can perform this action
	if !user.Authenticated() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// bind and validate review
	form := Review{}
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok, errors := ValidateReview(form)
	if !ok {
		c.HTML(http.StatusConflict, "review-form.html", gin.H{
			"Action": "Edit",
			"User":   user,
			"Form":   form,
			"Errors": errors,
		})
		return
	}

	// update review
	review, err := reviewStore.SaveReview(form)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	url := fmt.Sprintf("/whisky/%d", review.Whisky.ID)
	c.Redirect(http.StatusSeeOther, url)
}
