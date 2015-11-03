package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/auth0/go-jwt-middleware"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/dgrijalva/jwt-go"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	_ "github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/lib/pq"
)

var DB *sqlx.DB

func main() {
	port := os.Getenv("PORT")
	dburl := os.Getenv("DATABASE_URL")

	log.Printf("connecting to db: %s", dburl)
	db, err := sqlx.Connect("postgres", dburl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	cid := os.Getenv("AUTH0_CID")
	secret := os.Getenv("AUTH0_SECRET")
	auth := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			decoded, err := base64.URLEncoding.DecodeString(secret)
			if err != nil {
				return nil, err
			}
			return decoded, nil
		},
	})

	posts := PostCtrl{db}

	principal := PrincipalCtrl{db}
	friends := FriendsCtrl{db}

	users := UserCtrl{db}
	whiskeys := WhiskeyCtrl{db}

	r := gin.Default()
	r.Static("/app", "app")
	r.Static("/static", "static")
	r.StaticFile("/", "app/index.html")

	authorized := r.Group("/api")
	authorized.Use(Secured(cid, auth))
	{
		authorized.GET("/posts", posts.All)
		authorized.POST("/posts", posts.Create)
		authorized.DELETE("/posts/:id", posts.Delete)

		authorized.GET("/principal", principal.Get)
		authorized.PUT("/principal", principal.Create)

		authorized.GET("/principal/friends", friends.All)
		authorized.PUT("/principal/friends/:id", friends.Create)

		authorized.GET("/users", users.All)

		authorized.GET("/whiskeys", whiskeys.All)
		authorized.POST("/whiskeys", whiskeys.Create)
	}

	log.Printf("listening on :%s", port)
	r.Run(":" + port)
}

func Secured(aud string, auth *jwtmiddleware.JWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.CheckJWT(c.Writer, c.Request); err != nil {
			c.Abort()
			return
		}

		token := context.Get(c.Request, auth.Options.UserProperty).(*jwt.Token)
		if token.Claims["aud"].(string) != aud {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", token.Claims["sub"].(string))
		c.Next()
	}
}
