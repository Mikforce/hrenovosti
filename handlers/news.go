package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type News struct {
	ID       int    // Идентификатор новости
	Title    string // Заголовок новости
	ImageURL string // URL изображения новости
	Link     string // URL новости
	Source   string // Источник новости
}

func RandomNewsHandler(c *gin.Context) {
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	// Вызываем функцию getRandomNews
	result, err := GetRandomNews(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error getting random news: %v", err)})
		return
	}

	// Отправляем результат в шаблон

	c.HTML(http.StatusOK, "home.html", gin.H{"news": result})
}
func GetRandomNews(db *sql.DB) ([]News, error) {
	var id int
	var title, imageUrl, link string

	// Получаем три случайные записи из таблицы
	query := "SELECT id, title, image_url, link FROM news ORDER BY RANDOM() LIMIT 3"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении случайных записей: %v", err)
	}
	defer rows.Close()

	newsList := []News{}

	// Итерируем по результирующим строкам
	for rows.Next() {
		err := rows.Scan(&id, &title, &imageUrl, &link)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}

		// Формируем объект News и добавляем его в список
		news := News{Title: title, ImageURL: imageUrl, Link: link}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам результата: %v", err)
	}

	return newsList, nil
}

func GetAllArticlesHandler(c *gin.Context) {
	source := c.DefaultQuery("source", "")
	num := c.DefaultQuery("num", "10")
	numArticles, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'num' parameter"})
		return
	}

	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	var articles []News
	if source == "" {
		articles, err = GetAllArticles(db, numArticles)
	} else {
		articles, err = GetArticlesBySource(db, source, numArticles)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error getting articles: %v", err)})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func GetArticlesBySource(db *sql.DB, source string, numArticles int) ([]News, error) {
	query := "SELECT id, title, image_url, link, source FROM news WHERE source = ? LIMIT ?"

	rows, err := db.Query(query, source, numArticles)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	articlesList := []News{}
	for rows.Next() {
		var article News
		if err := rows.Scan(&article.ID, &article.Title, &article.ImageURL, &article.Link, &article.Source); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		articlesList = append(articlesList, article)
	}

	return articlesList, nil
}

func GetAllArticles(db *sql.DB, numArticles int) ([]News, error) {
	// Выполните запрос, чтобы получить все статьи из базы данных
	query := fmt.Sprintf("SELECT id, title, image_url, link FROM news LIMIT %d", numArticles)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	articlesList := []News{}

	// Перебирает строки результатов и создает объекты новостей.
	for rows.Next() {
		var id int
		var title, imageURL, link string

		err := rows.Scan(&id, &title, &imageURL, &link)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Create a News object and add it to the list
		article := News{Title: title, ImageURL: imageURL, Link: link}
		articlesList = append(articlesList, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through result rows: %v", err)
	}

	return articlesList, nil
}
