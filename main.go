package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

type News struct {
	Title    string
	ImageURL string
	Link     string
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
func getRandomNews(db *sql.DB) (News, error) {
	var id int
	var title, imageUrl, link string

	// Получаем количество записей в таблице
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
	if err != nil {
		return News{}, fmt.Errorf("ошибка при получении количества новостей: %v", err)
	}

	// Генерируем случайный индекс от 1 до count
	randomIndex := rand.Intn(count) + 1

	// Выполняем запрос на выборку рандомной записи
	query := fmt.Sprintf("SELECT id, title, image_url, link FROM news WHERE id = %d", randomIndex)
	err = db.QueryRow(query).Scan(&id, &title, &imageUrl, &link)
	if err != nil {
		return News{}, fmt.Errorf("ошибка при получении рандомной записи: %v", err)
	}

	// Формируем результат
	result := News{Title: title, ImageURL: imageUrl, Link: "https://panorama.pub" + link}
	return result, nil
}

func basa() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Определяем URL страницы, которую будем парсить
		url := "https://panorama.pub/science"

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

		// Устанавливаем соединение с базой данных
		db, err := sql.Open("sqlite3", "news.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// // Создаем таблицу news, если она еще не существует
		// _, err = db.Exec(`
		// 	CREATE TABLE IF NOT EXISTS news (
		// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
		// 		title TEXT,
		// 		image_url TEXT,
		// 		link TEXT
		// 	)
		// `)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		newsList := []News{}

		// Find the news items
		doc.Find(".grid-cols-1").Each(func(i int, s *goquery.Selection) {
			if i < 3 {
				// Для каждого найденного элемента получите заголовок, URL-адрес изображения и URL-адрес ссылки.
				title := s.Find("a").Text()
				imageURL, _ := s.Find("img").Attr("src")
				link, _ := s.Find("a").Attr("href")

				// Проверяем, существует ли уже запись с такими же данными
				var count int
				db.QueryRow("SELECT COUNT(*) FROM news WHERE title = ? AND image_url = ? AND link = ?", title, imageURL, link).Scan(&count)
				if count > 0 {
					// Запись уже существует, пропускаем ее
					return
				}

				// Выводим данные новости
				fmt.Println("Title:", title)
				fmt.Println("Image URL:", imageURL)
				fmt.Println("Link:", link)

				// Сохраняем данные в базу данных
				_, err = db.Exec(`
				INSERT INTO news (title, image_url, link)
				VALUES (?, ?, ?)
			`, title, imageURL, link)
				if err != nil {
					log.Fatal(err)

				}

				news := News{Title: title, ImageURL: imageURL, Link: "https://panorama.pub" + link}
				newsList = append(newsList, news)

			}
		})

		// Создаем шаблон HTML
		tmpl := template.Must(template.ParseFiles("html/home.html"))

		// Передаем список новостей в шаблон и генерируем HTML

		err = tmpl.Execute(w, newsList)
		if err != nil {
			log.Fatal(err)
		}
	})
}
func main() {

	http.HandleFunc("/", randomNewsHandler)
	// Устанавливаем путь к файлам шаблонов
	tmplPath := "html"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(tmplPath))))
	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
