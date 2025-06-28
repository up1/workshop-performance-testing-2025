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
)

var db *sql.DB

func main() {
	var err error
	dsn := os.Getenv("MYSQL_DSN") // e.g. user:pass@tcp(127.0.0.1:3306)/testdb
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	r := gin.Default()
	r.GET("/users", getUsersHandler)
	r.Run(":8080")
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUsersHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

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
