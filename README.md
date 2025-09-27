# 🧠 Auth Service  
_or: How I Learned to Stop Worrying and Love the Backend_

> "Your refresh token is a silly place!"  
> — Monty Python & the Holy Golang

## 🧾 TL;DR

Минималистичный сервис аутентификации, написанный на Go, потому что Python слишком мягкий, а C — это уже насилие.  

## 🧰 Технологии

| Штука        | Комментарий |
|--------------|-------------|
| Go           | Просто, быстро, неожиданно приятно |
| PostgreSQL   | Open source рулит|
| Docker       | Контейнеризация |
| JWT (SHA512) | аутентификация |
| bcrypt       | хлеб с солью |


## 🛠 Установка

```bash
docker-compose -f docker-compose.yml up -d
```

📦 Что умеет?

    🎫 Выдаёт пару токенов (access + refresh) по GUID'у

    🔄 Обновляет токены

    🔐 Защищает маршрут получения GUID — только с access токеном

    🗑 Удаляет сессию по access токену

🔐 Особенности токенов

    Access token:

        Формат — JWT
        Подпись — SHA512

    Refresh token:

        base64 снаружи, bcrypt внутри
        Защищен от повторного использования
        Хакерский подход с User-Agent в куках и IPшник палятся моментально
        Если IP странный — шлём webhook, но не паникуем)

📎 Почему Go?

Потому что я так думаю.

📃 Swagger

Чтобы была дока. 
