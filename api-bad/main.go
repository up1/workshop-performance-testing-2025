// main.go
package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zsais/go-gin-prometheus"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(db:3306)/demo")
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// NewWithConfig is the recommended way to initialize the middleware
	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	p.Use(r)

	r.GET("/users", getUsers)
	r.GET("/users/:id", getUserByID)
	r.Run(":8080")
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	var name string
	err := db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "query error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM users limit 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query error"})
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		users = append(users, map[string]interface{}{
			"id": id, "name": name,
		})
	}

	c.JSON(http.StatusOK, users)
}
