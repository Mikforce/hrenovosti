import requests
from bs4 import BeautifulSoup
import sqlite3

# URL страницы с новостями
url = 'https://ria.ru/world/'

response = requests.get(url)
html_content = response.text

soup = BeautifulSoup(html_content, 'html.parser')

news_items = soup.find_all(class_='list-item')

# Создаем или подключаемся к базе данных
conn = sqlite3.connect('news.db')
cursor = conn.cursor()

# Создаем таблицу, если она не существует
cursor.execute('''
    CREATE TABLE IF NOT EXISTS news (
        id INTEGER PRIMARY KEY,
        title TEXT,
        link TEXT UNIQUE,  -- Make 'link' column unique to prevent duplicates
        image_url TEXT
    )
''')

for news in news_items:
    title = news.find(class_='list-item__title').text.strip()
    link = news.find('a')['href']

    # Находим тег с изображением, если оно есть
    image_tag = news.find('img')
    if image_tag:
        image_url = image_tag['src']
    else:
        image_url = None

    # Проверяем, есть ли новость с такой ссылкой в базе данных
    existing_news = cursor.execute('SELECT id FROM news WHERE link = ?', (link,)).fetchone()

    if existing_news:
        print(f"Новость с ссылкой '{link}' уже существует. Пропускаем.")
    else:
        # Вставляем данные в базу данных
        cursor.execute('INSERT INTO news (title, link, image_url) VALUES (?, ?, ?)', (title, link, image_url))

        print(f"Заголовок: {title}")
        print(f"Ссылка: {link}")
        print(f"Изображение: {image_url}")
        print("=" * 50)

# Сохраняем изменения и закрываем соединение с базой данных
conn.commit()
conn.close()