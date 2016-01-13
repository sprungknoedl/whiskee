package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func PostgresDatabase() gin.HandlerFunc {
	url := os.Getenv("DATABASE_URL")
	store, err := NewPostgresStore(url)
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set(UserStoreKey, store)
		c.Set(WhiskyStoreKey, store)
		c.Set(ReviewStoreKey, store)
		c.Next()
	}
}

type PostgresStore struct {
	getUser    *sqlx.Stmt
	getAllUser *sqlx.Stmt
	insertUser *sqlx.NamedStmt

	getWhisky         *sqlx.Stmt
	getAllWhisky      *sqlx.Stmt
	getTrendingWhisky *sqlx.Stmt
	insertWhisky      *sqlx.NamedStmt
	updateWhisky      *sqlx.NamedStmt

	getReview           *sqlx.Stmt
	getAllReview        *sqlx.Stmt
	getActivity         *sqlx.Stmt
	insertReview        *sqlx.NamedStmt
	updateReview        *sqlx.NamedStmt
	updateReviewSummary *sqlx.Stmt
}

func NewPostgresStore(url string) (*PostgresStore, error) {
	store := &PostgresStore{}
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("connection: %v", err)
	}

	// statements for user store
	store.getUser, err = db.Preparex(`
		SELECT * 
		FROM users
		WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("sql get user: %v", err)
	}

	store.getAllUser, err = db.Preparex(`
		SELECT
			users.*,
			COUNT(reviews.*) as ratings,
			COUNT(NULLIF(reviews.description, '')) as reviews
		FROM users
		LEFT JOIN reviews ON users.id = reviews.user_id
		GROUP BY users.id
		ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("sql get all user: %v", err)
	}

	store.insertUser, err = db.PrepareNamed(`
		INSERT 
		INTO users (
			auth0,
			name,
			email,
			picture)
		VALUES (
			:auth0,
			:name,
			:email,
			:picture)
		ON CONFLICT (auth0) DO
		UPDATE SET
			name = :name,
			email = :email,
			picture = :picture
		RETURNING *`)
	if err != nil {
		return nil, fmt.Errorf("sql insert user: %v", err)
	}

	// statements for whisky store
	store.getWhisky, err = db.Preparex(`
		SELECT *
		FROM whiskies
		WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("sql get whisky: %v", err)
	}

	store.getAllWhisky, err = db.Preparex(`
		SELECT *
		FROM whiskies
		ORDER BY distillery, name, age`)
	if err != nil {
		return nil, fmt.Errorf("sql get all whisky: %v", err)
	}

	store.getTrendingWhisky, err = db.Preparex(`
		SELECT *
		FROM whiskies
		ORDER BY id DESC
		LIMIT $1`)
	if err != nil {
		return nil, fmt.Errorf("sql get trending whisky: %v", err)
	}

	store.insertWhisky, err = db.PrepareNamed(`
		INSERT 
		INTO whiskies (type, distillery, name, age, abv, description, picture, thumbnail)
	  VALUES (:type, :distillery, :name, :age, :abv, :description, :picture, :thumbnail)
		RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("sql insert whisky: %v", err)
	}

	store.updateWhisky, err = db.PrepareNamed(`
		UPDATE whiskies 
		SET type = :type, distillery = :distillery, name = :name,
				age = :age, abv = :abv, description = :description,
				picture = :picture, thumbnail = :thumbnail
		WHERE id = :id`)
	if err != nil {
		return nil, fmt.Errorf("sql update whisky: %v", err)
	}

	// statements for review store
	store.getReview, err = db.Preparex(`
		SELECT *
		FROM reviews
		WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("sql get review: %v", err)
	}

	store.getAllReview, err = db.Preparex(`
		SELECT *
		FROM reviews
		WHERE whisky_id = $1
		ORDER BY date DESC
		LIMIT $2`)
	if err != nil {
		return nil, fmt.Errorf("sql get all review: %v", err)
	}

	store.getActivity, err = db.Preparex(`
		SELECT *
		FROM reviews
		ORDER BY date DESC
		LIMIT $1`)
	if err != nil {
		return nil, fmt.Errorf("sql get all review: %v", err)
	}

	store.insertReview, err = db.PrepareNamed(`
		INSERT
		INTO reviews (user_id, whisky_id, rating, description)
		VALUES (:user_id, :whisky_id, :rating, :description)
		RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("sql insert review: %v", err)
	}

	store.updateReview, err = db.PrepareNamed(`
		UPDATE reviews
		SET rating = :rating, description = :description
		WHERE id = :id`)
	if err != nil {
		return nil, fmt.Errorf("sql update review: %v", err)
	}

	store.updateReviewSummary, err = db.Preparex(`
		WITH tmp AS (SELECT 
			 	COALESCE(AVG(rating), 0) as rating, 
			 	COUNT(*) as ratings,
				COUNT(NULLIF(description, '')) as reviews
			FROM reviews
			WHERE whisky_id = $1)
		UPDATE whiskies
		SET rating = tmp.rating, ratings = tmp.ratings, reviews = tmp.reviews
		FROM tmp
		WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("sql update review summary: %v", err)
	}

	return store, nil
}

// ---------------------------------------------------
// User Store
// ---------------------------------------------------
func (store PostgresStore) GetUser(id int) (User, error) {
	user := User{}
	err := store.getUser.Get(&user, id)
	return user, err
}

func (store PostgresStore) GetAllUser() ([]User, error) {
	users := []User{}
	err := store.getAllUser.Select(&users)
	return users, err
}

func (store PostgresStore) SaveUser(entity User) (User, error) {
	err := store.insertUser.Get(&entity, entity)
	return entity, err
}

// ---------------------------------------------------
// Whisky Store
// ---------------------------------------------------
func (store PostgresStore) GetWhisky(id int) (Whisky, error) {
	whisky := Whisky{}
	err := store.getWhisky.Get(&whisky, id)
	return whisky, err
}

func (store PostgresStore) GetAllWhiskies() ([]Whisky, error) {
	whiskies := []Whisky{}
	err := store.getAllWhisky.Select(&whiskies)
	return whiskies, err
}

func (store PostgresStore) GetTrending(limit int) ([]Whisky, error) {
	whiskies := []Whisky{}
	err := store.getTrendingWhisky.Select(&whiskies, limit)
	return whiskies, err
}

func (store PostgresStore) SaveWhisky(entity Whisky) (Whisky, error) {
	if entity.ID != 0 {
		_, err := store.updateWhisky.Exec(entity)
		return entity, err
	} else {
		err := store.insertWhisky.Get(&entity, entity)
		return entity, err
	}
}

// ---------------------------------------------------
// Review Store
// ---------------------------------------------------
func (store PostgresStore) joinReview(review Review) (Review, error) {
	review.Whisky = &Whisky{}
	if err := store.getWhisky.Get(review.Whisky, review.WhiskyID); err != nil {
		return review, err
	}

	review.User = &User{}
	if err := store.getUser.Get(review.User, review.UserID); err != nil {
		return review, err
	}

	return review, nil
}

func (store PostgresStore) GetReview(id int) (Review, error) {
	review := Review{}
	err := store.getReview.Get(&review, id)
	if err != nil {
		return review, err
	}

	return store.joinReview(review)
}

func (store PostgresStore) GetAllReviews(whiskyid int, limit int) ([]Review, error) {
	reviews := []Review{}
	err := store.getAllReview.Select(&reviews, whiskyid, limit)
	if err != nil {
		return reviews, err
	}

	for i, review := range reviews {
		reviews[i], err = store.joinReview(review)
		if err != nil {
			return reviews, err
		}
	}

	return reviews, nil
}

func (store PostgresStore) GetActivity(userid int, limit int) ([]Review, error) {
	reviews := []Review{}
	err := store.getActivity.Select(&reviews, limit)
	if err != nil {
		return reviews, err
	}

	for i, review := range reviews {
		reviews[i], err = store.joinReview(review)
		if err != nil {
			return reviews, err
		}
	}

	return reviews, nil
}

func (store PostgresStore) SaveReview(entity Review) (Review, error) {
	if entity.ID != 0 {
		if _, err := store.updateReview.Exec(entity); err != nil {
			return entity, err
		}
	} else {
		if err := store.insertReview.Get(&entity, entity); err != nil {
			return entity, err
		}
	}

	if _, err := store.updateReviewSummary.Exec(entity.WhiskyID); err != nil {
		return entity, err
	}

	return store.joinReview(entity)
}

var (
	_ UserStore   = (*PostgresStore)(nil)
	_ WhiskyStore = (*PostgresStore)(nil)
	_ ReviewStore = (*PostgresStore)(nil)
)
