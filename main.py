import requests
from bs4 import BeautifulSoup
import sqlite3

# Определяем URL страницы, которую будем парсить
url = "https://ria.ru/organization_API"

# Отправляем GET-запрос и получаем содержимое страницы
response = requests.get(url)
page_content = response.content

# Создаем объект BeautifulSoup для парсинга HTML-кода
soup = BeautifulSoup(page_content, "html.parser")

# Находим все элементы с классом "list-item"
news_items = soup.find_all(class_="list-item")

# Устанавливаем соединение с базой данных
conn = sqlite3.connect("news.db")
cursor = conn.cursor()

# Создаем таблицу news, если она еще не существует
cursor.execute("""
    CREATE TABLE IF NOT EXISTS news (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        image_url TEXT,
        link TEXT
    )
""")

# Проходимся по каждому элементу "list-item" и извлекаем необходимые данные
for item in news_items:
    # Извлекаем заголовок новости
    title = item.find(class_="list-item__title").text.strip()

    # Извлекаем ссылку на новость
    link = item.find("a")["href"]

    # Извлекаем URL изображения
    image_url = item.find("img")["src"]

    # Вставляем данные о новости в таблицу
    cursor.execute("""
        INSERT INTO news (title, image_url, link)
        VALUES (?, ?, ?)
    """, (title, image_url, link))

    # Сохраняем изменения
    conn.commit()

# Закрываем соединение с базой данных
conn.close()