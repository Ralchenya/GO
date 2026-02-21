package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Важно: импортируем с подчеркиванием
)

func main() {
	// 1. Подключаемся к базе данных (файл test.db будет создан автоматически)
	// Важно: Для SQLite рекомендуется ограничивать количество соединений
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	// Закрываем соединение, когда функция main завершится
	defer db.Close()

    // Устанавливаем лимит соединений (для SQLite рекомендуется 1)
    db.SetMaxOpenConns(1)

	// 2. Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	fmt.Println("Успешно подключились к SQLite!")

	// 3. Создаем таблицу
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
	fmt.Println("Таблица users создана или уже существует.")

	// 4. Вставляем данные
	insertUserSQL := `INSERT INTO users (name, email) VALUES (?, ?)`
	result, err := db.Exec(insertUserSQL, "Тестовый Пользователь", "test@example.com")
	if err != nil {
		log.Fatal("Ошибка вставки данных:", err)
	}

	userID, _ := result.LastInsertId()
	fmt.Printf("Добавлен пользователь с ID: %d\n", userID)

	// 5. Читаем данные
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal("Ошибка чтения данных:", err)
	}
	defer rows.Close()

	fmt.Println("Список пользователей:")
	for rows.Next() {
		var id int
		var name string
		var email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal("Ошибка сканирования строки:", err)
		}
		fmt.Printf("  %d: %s (%s)\n", id, name, email)
	}
}