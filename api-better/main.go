package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/zsais/go-gin-prometheus"
)

var db *sql.DB
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // e.g. "localhost:6379"
		DB:   0,
	})
	// Test the Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Cannot connect to Redis: ", err)
	}
}

func main() {
	var err error
	dsn := os.Getenv("MYSQL_DSN") // e.g. user:pass@tcp(127.0.0.1:3306)/testdb
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	initRedis()

	// Set connection pool parameters
	// Adjust these values based on your application's needs
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	// This is important to ensure the database is reachable before starting the server
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	// Create a new Gin router
	// Use gin.New() to create a new router without default middleware
	r := gin.New()
	r.Use(gin.Recovery())

	// NewWithConfig is the recommended way to initialize the middleware
	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	p.Use(r)

	r.GET("/users", getUsersHandler)
	r.Run(":8080")
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUsersHandler(c *gin.Context) {
	// Use a context with timeout for the database query
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	// 1. Try Redis first
	const cacheKey = "users:all"
	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	// 2. Fallback to DB
	users, err := fetchUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func fetchUsers(ctx context.Context) ([]User, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
