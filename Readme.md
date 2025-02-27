# **Lock-Stock v2**

## **О проекте**

Lock-Stock v2 — это backend-приложение для игры, основанной на популярном YouTube-шоу "Lock stock". Основная цель
проекта — предоставить игровую механику, управляемую через API.


---

## **Технологии**

- **Язык:** Go
- **База данных:** PostgreSQL
- **Контейнеризация:** Docker + Docker Compose
- **Миграции:** golang-migrate
- **Сетевой фреймворк:** Chi
- **Dependency Injection:** Wire

---

### API маршруты

Проект включает следующие основные маршруты:

- [WebSocket маршруты](app/handlers/http/ws/ws.yaml)
- [Маршруты пользователя](app/handlers/http/user/user.yaml)
- [Маршруты комнат](app/handlers/http/room/room.yaml)

---

## **Запуск проекта**

### **1. Установка**

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/username/lock-stock-v2.git
   ```
2. Запустите команду:
   ```bash
   make init
   ```
   
---

## **Планы на будущее**
1. **Добавить тесты:**
   - Покрыть unit-тестами основные use-case сервисы.
   - Добавить интеграционные тесты.

2. **Расширить функционал:**
   - Добавить новые API-эндпоинты.
   - Перенести отправку сообщений в web socket на kafka? 

3. **Документация:**
   - ✔ ~~Добавить подробное описание всех маршрутов.~~
   - ✔ ~~Генерация кода на основе yaml файлов~~

4. **Non-dev environment**
   - Настройка docker окружения для qa и prod окружений
   - Настройка подключения к бд для qa и prod окружений
   - Настройка подключения зависимостей для каждого окружения

---

## **Контакты**
Если у вас есть вопросы или предложения, свяжитесь с разработчиком:

- GitHub: [scarymovie](https://github.com/scarymovie)
- Email: vino.zeka@gmail.com

