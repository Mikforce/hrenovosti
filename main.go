package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/PuerkitoBio/goquery"
// )

// func ExampleScrape(w http.ResponseWriter, r *http.Request) {
// 	// Request the HTML page.
// 	res, err := http.Get("https://panorama.pub")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != 200 {
// 		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
// 	}

// 	// Load the HTML document
// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create an HTML string to store the parsed data
// 	html := "<html><body>"

// 	// Find the news items
// 	doc.Find(".grid-cols-1").Each(func(i int, s *goquery.Selection) {
// 		// For each item found, get the title and image URL
// 		title := s.Find("a").Text()
// 		imageURL, _ := s.Find("img").Attr("src")

// 		// Append the title and image URL to the HTML string
// 		html += fmt.Sprintf("<h2>%s</h2>", title)
// 		html += fmt.Sprintf(`<img src="%s" alt="%s">`, imageURL, title)
// 	})

// 	html += "</body></html>"

// 	// Set the response headers
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	// Write the HTML string to the response
// 	w.Write([]byte(html))
// }

// func main() {
// 	http.HandleFunc("/", ExampleScrape)
// 	log.Fatal(http.ListenAndServe(":80", nil))
// }

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Определяем URL страницы, которую будем парсить
	url := "https://ria.ru/organization_API"

	// Отправляем GET-запрос и получаем содержимое страницы
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Создаем объект goquery для парсинга HTML-кода
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Находим все элементы с классом "list-item"
	newsItems := doc.Find(".list-item")

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("sqlite3", "news.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу news, если она еще не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS news (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			image_url TEXT,
			link TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Проходимся по каждому элементу "list-item" и извлекаем необходимые данные
	newsItems.Each(func(i int, item *goquery.Selection) {
		// Извлекаем заголовок новости
		title := item.Find(".list-item__title").Text()

		// Извлекаем ссылку на новость
		link, _ := item.Find("a").Attr("href")

		// Извлекаем URL изображения
		imageURL, _ := item.Find("img").Attr("src")

		// Вставляем данные о новости в таблицу
		_, err := db.Exec(`
			INSERT INTO news (title, image_url, link)
			VALUES (?, ?, ?)
		`, title, imageURL, link)
		if err != nil {
			log.Fatal(err)
		}
	})

	fmt.Println("Parsing completed successfully.")
}
