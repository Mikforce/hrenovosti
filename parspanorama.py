import requests
from bs4 import BeautifulSoup
import sqlite3

# Создаем базу данных или подключаемся к существующей
conn = sqlite3.connect('news.db')
cursor = conn.cursor()

# Создаем таблицу для новостей, если она еще не существует
cursor.execute('''
    CREATE TABLE IF NOT EXISTS news (
        id INTEGER PRIMARY KEY,
        title TEXT,
        image_url TEXT,
        link TEXT UNIQUE
    )
''')
conn.commit()

base_url = 'https://panorama.pub'
page_num = 2  # Начнем с 2 страницы согласно вашему примеру
url = f'{base_url}/science?page={page_num}'
response = requests.get(url)
soup = BeautifulSoup(response.content, 'html.parser')

news_list = soup.find_all('a', class_='hover:text-secondary')

for news in news_list:
    title = news.find('div', class_='text-xl').text.strip()
    image_url = news.find('img')['src']
    link = news['href']

    try:
        cursor.execute('INSERT INTO news (title, image_url, link) VALUES (?, ?, ?)', (title, image_url, "https://panorama.pub" + link))
        conn.commit()
    except sqlite3.IntegrityError:
        print(f"Новость с ссылкой '{link}' уже существует в базе данных. Пропускаю...")

conn.close()