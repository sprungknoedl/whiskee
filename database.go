package main

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(User{})
}

/****************************************************
* DATABASE                                          *
****************************************************/
type User struct {
	ID      int    `db:"id"`
	Auth0   string `db:"auth0"`
	Name    string `db:"name"`
	Email   string `db:"email"`
	Picture string `db:"picture"`
}

func (u User) Authenticated() bool {
	return u.ID != 0
}

type Whisky struct {
	ID          int     `db:"id" form:"id"`
	Type        string  `db:"type" form:"type"`
	Distillery  string  `db:"distillery" form:"distillery"`
	Name        string  `db:"name" form:"name"`
	Age         int     `db:"age" form:"age"`
	ABV         float64 `db:"abv" form:"abv"`
	Description string  `db:"description" form:"description"`

	Picture   string `db:"picture" form:"picture"`
	Thumbnail string `db:"thumbnail" form:"thumbnail"`

	Rating  float64 `db:"rating"`
	Ratings int     `db:"ratings"`
	Reviews int     `db:"reviews"`
}

type Review struct {
	ID          int       `db:"id" form:"id"`
	Date        time.Time `db:"date"`
	UserID      int       `db:"user_id"`
	WhiskyID    int       `db:"whisky_id" form:"whisky"`
	Rating      int       `db:"rating" form:"rating"`
	Description string    `db:"description" form:"description"`

	User   *User
	Whisky *Whisky
}

var (
	UserStoreKey   = "github.com/sprungknoedl/whiskee/user-store"
	WhiskyStoreKey = "github.com/sprungknoedl/whiskee/whisky-store"
	ReviewStoreKey = "github.com/sprungknoedl/whiskee/review-store"
)

type UserStore interface {
	GetUser(int) (User, error)
	SaveUser(User) (User, error)
}

type WhiskyStore interface {
	GetWhisky(int) (Whisky, error)
	GetAllWhiskies() ([]Whisky, error)
	SaveWhisky(Whisky) (Whisky, error)

	GetTrending(limit int) ([]Whisky, error)
}

type ReviewStore interface {
	GetReview(int) (Review, error)
	SaveReview(Review) (Review, error)

	GetAllReviews(whiskyid int, limit int) ([]Review, error)
	GetActivity(userid int, limit int) ([]Review, error)
}

/****************************************************
* VALIDATION                                        *
****************************************************/

var (
	ValidationRequired    = "cannot be empty"
	ValidationNonZero     = "cannot be zero"
	ValidationNonNegative = "must be greater than zero"
)

func ValidateWhisky(form Whisky) (bool, map[string]string) {
	errors := map[string]string{}

	if form.Type == "" {
		errors["type"] = ValidationRequired
	}

	if form.Distillery == "" {
		errors["distillery"] = ValidationRequired
	}

	if form.Age == 0 {
		errors["age"] = ValidationNonZero
	} else if form.Age < 0 {
		errors["age"] = ValidationNonNegative
	}

	if form.ABV < 0 {
		errors["abv"] = ValidationNonNegative
	}

	return len(errors) == 0, errors
}

func ValidateReview(form Review) (bool, map[string]string) {
	errors := map[string]string{}
	return len(errors) == 0, errors
}
