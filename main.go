package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type News struct {
	Title    string
	ImageURL string
	Link     string
}

func parseNews(url string) ([]News, error) {
	// Отправляем GET-запрос и получаем содержимое страницы
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Создаем объект goquery для парсинга HTML-кода
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	newsList := []News{}

	// Найти новости
	doc.Find(".flex-col").Each(func(i int, s *goquery.Selection) {
		if i < 3 {
			// Для каждого найденного элемента получите заголовок, URL-адрес изображения и URL-адрес ссылки.
			title := s.Find(".text-xl").Text()
			imageURL, _ := s.Find("img").Attr("src")
			link, _ := s.Find("a").Attr("href")

			news := News{Title: title, ImageURL: imageURL, Link: link}
			newsList = append(newsList, news)
		}
	})

	return newsList, nil
}

func baze() {
	baseURL := "https://panorama.pub/science?page="
	rand.Seed(time.Now().UnixNano())
	numPages := rand.Intn(30) + 1

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблици вставляем данные
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS news ( 		id INTEGER PRIMARY KEY AUTOINCREMENT, 		title TEXT, 		image_url TEXT, 		link TEXT 	)")
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= numPages; i++ {
		pageURL := fmt.Sprintf("%s%d", baseURL, i)
		newsList, err := parseNews(pageURL)
		if err != nil {
			log.Println("Error parsing news:", err)
			continue
		}

		// Вставить данные в базу данных
		for _, news := range newsList {
			if news.Title != "" { // Проверяем, не является ли заголовок новости пустым
				_, err = db.Exec("INSERT INTO news (title, image_url, link) VALUES (?, ?, ?)", news.Title, news.ImageURL, news.Link)
				if err != nil {
					log.Println("Error inserting data into database:", err)
				}
			}
		}
	}
}
func randomNewsHandler(c *gin.Context) {
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	// Вызываем функцию getRandomNews
	result, err := getRandomNews(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error getting random news: %v", err)})
		return
	}

	// Отправляем результат в шаблон

	c.HTML(http.StatusOK, "home.html", gin.H{"news": result})
}
func getRandomNews(db *sql.DB) ([]News, error) {
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
		news := News{Title: title, ImageURL: imageUrl, Link: "https://panorama.pub" + link}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам результата: %v", err)
	}

	return newsList, nil
}

func getAllArticlesHandler(c *gin.Context) {
	// Получить параметр 'num' из строки запроса
	num := c.DefaultQuery("num", "3") // По умолчанию 3, если «число» не указано
	// Преобразование «число» в целое число
	numArticles, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'num' parameter"})
		return
	}

	// Подключиться к базе данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	// Get all articles from the database
	articles, err := getAllArticles(db, numArticles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error getting articles from database: %v", err)})
		return
	}

	// Send the list of articles in JSON format
	c.JSON(http.StatusOK, articles)
}
func getAllArticles(db *sql.DB, numArticles int) ([]News, error) {
	// Execute the query to get all articles from the database
	query := fmt.Sprintf("SELECT id, title, image_url, link FROM news LIMIT %d", numArticles)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	articlesList := []News{}

	// Iterate through the result rows and create news objects.
	for rows.Next() {
		var id int
		var title, imageURL, link string

		err := rows.Scan(&id, &title, &imageURL, &link)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Create a News object and add it to the list
		article := News{Title: title, ImageURL: imageURL, Link: "https://panorama.pub" + link}
		articlesList = append(articlesList, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through result rows: %v", err)
	}

	return articlesList, nil
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*") // Load HTML templates

	go baze()
	r.DELETE("/articles/delete", deleteArticleHandler)

	r.GET("/", randomNewsHandler)
	r.GET("/articles", getAllArticlesHandler)
	r.Static("/static", "./static") // Serve static files

	// New endpoint to add an article
	r.POST("/articles/add", addArticleHandler)

	fmt.Println("Server started on http://localhost:80")
	log.Fatal(r.Run(":80"))
}

type Article struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Link     string `json:"link"`
}

func addArticleHandler(c *gin.Context) {
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
func deleteArticleHandler(c *gin.Context) {
	// Get the article ID from the query parameter
	articleID := c.Query("id")

	// Connect to the database
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error connecting to database: %v", err)})
		return
	}
	defer db.Close()

	// Delete the article from the database
	_, err = db.Exec("DELETE FROM news WHERE id = ?", articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error deleting article from database: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// )

// func main() {
// 	// URL страницы с новостями
// 	url := "https://ria.ru/world/"

// 	// Получаем HTML-контент страницы
// 	doc, err := goquery.NewDocument(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Создаем или подключаемся к базе данных SQLite3
// 	db, err := sql.Open("sqlite3-lite", "news.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Создаем таблицу, если она не существует
// 	_, err = db.Exec(`
// 		CREATE TABLE IF NOT EXISTS news1 (
// 			id INTEGER PRIMARY KEY,
// 			title TEXT,
// 			link TEXT,
// 			image_url TEXT
// 		)
// 	`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Итерируемся по элементам с классом "list-item"
// 	doc.Find(".list-item").Each(func(i int, s *goquery.Selection) {
// 		// Получаем заголовок и ссылку
// 		title := strings.TrimSpace(s.Find(".list-item__title").Text())
// 		link, _ := s.Find("a").Attr("href")

// 		// Получаем ссылку на изображение, если оно есть
// 		imageURL, _ := s.Find("img").Attr("src")

// 		// Вставляем данные в базу данных
// 		_, err := db.Exec("INSERT INTO news (title, link, image_url) VALUES (?, ?, ?)", title, link, imageURL)
// 		if err != nil {
// 			log.Println("Ошибка при вставке данных:", err)
// 		}

// 		// Выводим информацию
// 		fmt.Println("Заголовок:", title)
// 		fmt.Println("Ссылка:", link)
// 		fmt.Println("Изображение:", imageURL)
// 		fmt.Println("=")
// 	})
// }
