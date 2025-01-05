package upload

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/defskela/httpServer/router"
	httpServer "github.com/defskela/httpServer/server"
)

func TestFileUpload(t *testing.T) {
	log.Println("Инициализация маршрутизатора и сервера...")
	router := router.NewRouter()
	router.Post("/upload", UploadHandler)
	server := httpServer.NewServer(router)

	log.Println("Запуск сервера на порту 8080...")
	server.Start("8080")
	defer func() {
		log.Println("Завершение работы сервера...")
		time.Sleep(1 * time.Second) // Задержка перед завершением
		server.Shutdown()
	}()

	filePath := "test.txt"
	fileContent := []byte("This is a test file content")
	log.Printf("Создание тестового файла: %s\n", filePath)
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		t.Fatalf("Ошибка создания тестового файла: %v", err)
	}

	defer func() {
		log.Printf("Удаление тестового файла: %s\n", filePath)
		err := os.Remove(filePath)
		if err != nil {
			t.Fatalf("Ошибка удаления тестового файла: %v", err)
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	log.Println("Создание формы для отправки файла...")
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		t.Fatalf("Ошибка создания части формы: %v", err)
	}

	log.Printf("Открытие файла для чтения: %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Ошибка открытия тестового файла: %v", err)
	}
	defer func() {
		log.Printf("Закрытие файла: %s\n", filePath)
		file.Close()
	}()

	log.Println("Копирование содержимого файла в форму...")
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Ошибка копирования данных файла: %v", err)
	}

	log.Println("Закрытие writer для завершения формирования запроса...")
	err = writer.Close()
	if err != nil {
		t.Fatalf("Ошибка закрытия writer: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
	if err != nil {
		t.Fatalf("Ошибка создания HTTP-запроса: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Println("Отправка HTTP-запроса на сервер...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка отправки HTTP-запроса: %v", err)
	}
	defer func() {
		log.Println("Закрытие ответа сервера...")
		resp.Body.Close()
	}()

	log.Printf("Код ответа сервера: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус код 200, получен: %d", resp.StatusCode)
	}

	log.Println("Тест upload успешно пройден")
}
