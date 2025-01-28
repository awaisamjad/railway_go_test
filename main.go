package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDatabase() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}
func main() {
	// Replace with your actual Railway database credentials
	dsn := "awaisamjad:Gunner$123@tcp(mariadb.railway.internal:3306)/whisp_db"
	var err error

	// Open the database connection
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Connected to MariaDB!")
	
	err = initDatabase()
	if err != nil {
		log.Fatalf("Error init database : %v", err)
	}

	// Initialize the Gin router
	router := gin.Default()

	// Route to display all posts
	router.GET("/", func(c *gin.Context) {
		posts, err := fetchPosts()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error fetching posts: %v", err)
			return
		}
		c.JSON(http.StatusOK, posts)
	})

	// Route to add a new post
	router.GET("/add-post", func(c *gin.Context) {
		err := addPost("This is a new post!")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error adding post: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	})

	// Start the server
	router.Run(":8080")
	log.Println("Running on http://localhost:8080")
}

// Function to fetch posts from the database
func fetchPosts() ([]string, error) {
	rows, err := db.Query("SELECT content FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		posts = append(posts, content)
	}
	return posts, nil
}

// Function to add a new post to the database
func addPost(content string) error {
	_, err := db.Exec("INSERT INTO posts (content) VALUES (?)", content)
	return err
}

