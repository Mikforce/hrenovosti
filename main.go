package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
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
func randomNewsHandler(w http.ResponseWriter, r *http.Request) {

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Вызываем функцию getRandomNews
	result, err := getRandomNews(db)
	if err != nil {
		log.Fatal(err)
	}

	// Отправляем результат в шаблон
	tmpl := template.Must(template.ParseFiles("html/home.html"))
	err = tmpl.Execute(w, result)
	if err != nil {
		log.Fatal(err)
	}
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

func main() {
	baze()
	http.HandleFunc("/", randomNewsHandler)
	// Устанавливаем путь к файлам шаблонов
	tmplPath := "html"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(tmplPath))))
	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
