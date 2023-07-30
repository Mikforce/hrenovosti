package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

type News struct {
	Title    string
	ImageURL string
	Link     string
}

func main() {
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
		if err != nil {
			log.Fatal(err)
		}

		newsList := []News{}

		// Find the news items
		doc.Find(".grid-cols-1").Each(func(i int, s *goquery.Selection) {
			if i < 3 {
				// For each item found, get the title, image URL, and link URL
				title := s.Find("a").Text()
				imageURL, _ := s.Find("img").Attr("src")
				link, _ := s.Find("a").Attr("href")

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

	fmt.Println("Server started on http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
