package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func gopython() {
	// Specify the Python script file to execute
	pythonScript := "parse_python/parsria.py"

	// Prepare the command to run the Python script
	cmd := exec.Command("/usr/bin/python3", pythonScript)

	// Set up pipes for standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the Python script
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
func gopythontwo() {
	// Specify the Python script file to execute
	pythonScript := "parse_python/parspanorama.py"

	// Prepare the command to run the Python script
	cmd := exec.Command("/usr/bin/python3", pythonScript)

	// Set up pipes for standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the Python script
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// func parseNews(url string) ([]News, error) {
// 	// Отправляем GET-запрос и получаем содержимое страницы
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	// Создаем объект goquery для парсинга HTML-кода
// 	doc, err := goquery.NewDocumentFromReader(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	newsList := []News{}

// 	// Найти новости
// 	doc.Find(".flex-col").Each(func(i int, s *goquery.Selection) {
// 		if i < 3 {
// 			// Для каждого найденного элемента получите заголовок, URL-адрес изображения и URL-адрес ссылки.
// 			title := s.Find(".text-xl").Text()
// 			imageURL, _ := s.Find("img").Attr("src")
// 			link, _ := s.Find("a").Attr("href")

// 			news := News{Title: title, ImageURL: imageURL, Link: link}
// 			newsList = append(newsList, news)
// 		}
// 	})

// 	return newsList, nil
// }

// func baze() {
// 	baseURL := "https://panorama.pub/science?page="
// 	rand.Seed(time.Now().UnixNano())
// 	numPages := rand.Intn(30) + 1

// 	// Устанавливаем соединение с базой данных
// 	db, err := sql.Open("sqlite3", "news.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Создаем таблици вставляем данные
// 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS news ( 		id INTEGER PRIMARY KEY AUTOINCREMENT, 		title TEXT, 		image_url TEXT, 		link TEXT 	)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for i := 1; i <= numPages; i++ {
// 		pageURL := fmt.Sprintf("%s%d", baseURL, i)
// 		newsList, err := parseNews(pageURL)
// 		if err != nil {
// 			log.Println("Error parsing news:", err)
// 			continue
// 		}

//			// Вставить данные в базу данных
//			for _, news := range newsList {
//				if news.Title != "" { // Проверяем, не является ли заголовок новости пустым
//					_, err = db.Exec("INSERT INTO news (title, image_url, link) VALUES (?, ?, ?)", news.Title, news.ImageURL, news.Link)
//					if err != nil {
//						log.Println("Error inserting data into database:", err)
//					}
//				}
//			}
//		}
//	}
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
		news := News{Title: title, ImageURL: imageUrl, Link: link}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам результата: %v", err)
	}

	return newsList, nil
}

func getAllArticlesHandler(c *gin.Context) {
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
		articles, err = getAllArticles(db, numArticles)
	} else {
		articles, err = getArticlesBySource(db, source, numArticles)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error getting articles: %v", err)})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func getArticlesBySource(db *sql.DB, source string, numArticles int) ([]News, error) {
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

func getAllArticles(db *sql.DB, numArticles int) ([]News, error) {
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

func main() {

	// go baze()
	go gopython()
	go gopythontwo()

	r := gin.Default()
	r.LoadHTMLGlob("html/*.html") // Load HTML templates

	r.DELETE("/articles/delete", deleteArticleHandler)

	r.GET("/", randomNewsHandler)
	r.GET("/articles", getAllArticlesHandler)
	r.Static("/static", "./static") // Serve static files

	// New endpoint to add an article
	r.POST("/articles/add", addArticleHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
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

// GET http://localhost:8080/articles
// Получение новостей по источнику:
// Чтобы получить новости по определенному источнику, отправьте GET-запрос на URL вашего веб-сервера с параметром source. Например:

// GET http://localhost:8080/articles?source=example_source
// Получение определенного количества новостей:
// Если вы хотите получить определенное количество новостей, отправьте GET-запрос на URL вашего веб-сервера с параметром num. Например, чтобы получить 5 новостей:

// GET http://localhost:8080/articles?num=5
// Получение новостей по источнику и количеству:
// Вы также можете комбинировать параметры source и num, чтобы получить новости по определенному источнику и указанному количеству. Например:

// GET http://localhost:8080/articles?source=example_source&num=3
