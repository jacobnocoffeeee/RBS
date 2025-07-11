package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

func connectDBWithRetry() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_NAME"),
		os.Getenv("PGSSLMODE"),
	)

	var db *sql.DB
	var err error
	for i := 1; i <= 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			return db, nil
		}
		log.Printf("Попытка подключения %d/10 не удалась: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
}

func main() {
	db, err := connectDBWithRetry()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, "<html><body><h1>Pong</h1></body></html>")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"checked": time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/list", func(c *gin.Context) {
		rows, err := db.Query("SELECT city, temperature FROM weather")
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
			return
		}
		defer rows.Close()

		html := "<html><body><h1>City --- Weather</h1><ul>"
		for rows.Next() {
			var city string
			var temp int
			if err := rows.Scan(&city, &temp); err != nil {
				log.Println("Ошибка чтения строки:", err)
				continue
			}
			html += fmt.Sprintf("<li>%s --- %d</li>", city, temp)
		}
		html += "</ul></body></html>"

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	// Простой POST запрос для добавления данных
	r.POST("/add", func(c *gin.Context) {
		city := c.PostForm("city")
		temp := c.PostForm("temp")

		_, err := db.Exec("INSERT INTO weather (city, temperature) VALUES ($1, $2)", city, temp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить данные"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Данные добавлены"})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}