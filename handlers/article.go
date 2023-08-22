package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Link     string `json:"link"`
}

func AddArticleHandler(c *gin.Context) {
	// Parse the JSON request body into the Article struct
	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	// Insert the article into the database
	_, err = db.Exec("INSERT INTO news (title, image_url, link) VALUES (?, ?, ?)", article.Title, article.ImageURL, article.Link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error inserting data into database: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article added successfully"})
}

func DeleteArticleHandler(c *gin.Context) {
	// Get the article ID from the query parameter

	intID := c.Query("id")

	// Connect to the database
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	// Delete the article from the database
	_, err = db.Exec("DELETE FROM news WHERE id = ?", intID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error deleting article from database: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

// func DeleteNewsHandler(c *gin.Context) {
// 	ID := c.Param("id") // Assuming the ID is passed as a URL parameter

// 	// Convert the ID string to an integer
// 	intID, err := strconv.Atoi(ID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}
// 	db := c.MustGet("db").(*sql.DB)
// 	// Call DelNews function passing the correct data type for ID
// 	_, err = DelNews(db, intID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
// }

// type Now struct {
// 	ID       int    // Идентификатор новости
// 	Title    string // Заголовок новости
// 	ImageURL string // URL изображения новости
// 	Link     string // URL новости
// 	Source   string // Источник новости
// }

// func DelNews(db *sql.DB, ID int) ([]Now, error) {
// 	// Выполните запрос, чтобы получить все статьи из базы данных
// 	query := fmt.Sprintf("SELECT id FROM news LIMIT %d", ID)

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("error executing query: %v", err)
// 	}
// 	defer rows.Close()

// 	articlesList := []Now{}

// 	// Перебирает строки результатов и создает объекты новостей.
// 	for rows.Next() {
// 		var id int
// 		var title, imageURL, link string

// 		err := rows.Scan(&id, &title, &imageURL, &link)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning row: %v", err)
// 		}

// 		// Удаление статьи из базы данных
// 		_, err = db.Exec("DELETE FROM news WHERE id = ?", id)
// 		if err != nil {
// 			return nil, fmt.Errorf("error deleting article from database: %v", err)
// 		}

// 	}

// 	return articlesList, nil
// }
