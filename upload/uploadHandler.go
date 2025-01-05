package upload

import (
	"filesStorage/utils"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/defskela/httpServer/logger"
	servModels "github.com/defskela/httpServer/models"
	servUtils "github.com/defskela/httpServer/utils"
)

func UploadHandler(conn net.Conn, data servModels.RequestData) {
	logger.Debug("DATA", data)
	conn.Write([]byte(servUtils.CreateHTTPResponse(200, "Information received")))

	fileName, exists := data.FormData["fileName"]
	if !exists {
		logger.Warn("Не передано имя файла")
		conn.Write([]byte(servUtils.CreateHTTPResponse(400, "No file name")))
		return
	}

	currentTime := time.Now().Format("2006_01_02_15_04_05")

	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		logger.Warn("Ошибка определения корня проекта:", err)
		conn.Write([]byte(servUtils.CreateHTTPResponse(500, "Error determining project root")))
		return
	}

	defaultDir := filepath.Join(projectRoot, "files")

	// Создание папки, если её нет
	err = os.MkdirAll(defaultDir, os.ModePerm)
	if err != nil {
		logger.Warn("Ошибка создания директории:", err)
		conn.Write([]byte(servUtils.CreateHTTPResponse(500, "Error creating directory")))
		return
	}

	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	fullFileName := filepath.Join(defaultDir, name+"_"+currentTime+ext)
	logger.Debug("Полное имя файла:", fullFileName)

	file, err := os.Create(fullFileName)
	if err != nil {
		logger.Warn("Ошибка создания файла:", err)
		conn.Write([]byte(servUtils.CreateHTTPResponse(500, "Error creating file")))
		return
	}

	defer file.Close()

	text, exists := data.FormData["file"]

	if !exists {
		logger.Warn("Не передано содержимое файла")
		conn.Write([]byte(servUtils.CreateHTTPResponse(400, "No file content")))
		return
	}

	_, err = file.WriteString(text)
	if err != nil {
		logger.Warn("Ошибка записи в файл:", err)
		conn.Write([]byte(servUtils.CreateHTTPResponse(500, "Error writing to file")))
		return
	}

	logger.Debug("Файл успешно создан и записан:", fileName)
}
