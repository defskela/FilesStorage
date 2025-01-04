package upload

import (
	"fmt"
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
		logger.Debug("Не передано имя файла")
		conn.Write([]byte(servUtils.CreateHTTPResponse(400, "No file name")))
		return
	}

	currentTime := time.Now().Format("2006_01_02_15_04_05")

	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	fullFileName := "files/" + name + "_" + currentTime + ext
	file, err := os.Create(fullFileName)
	if err != nil {
		logger.Debug("Ошибка создания файла:", err)
		conn.Write([]byte(servUtils.CreateHTTPResponse(500, "Error creating file")))
		return
	}

	defer file.Close()

	text, exists := data.FormData["file"]

	if !exists {
		logger.Debug("Не передано содержимое файла")
		conn.Write([]byte(servUtils.CreateHTTPResponse(400, "No file content")))
		return
	}

	_, err = file.WriteString(text)
	if err != nil {
		logger.Debug("Ошибка записи в файл:", err)
		conn.Write([]byte(servUtils.CreateHTTPResponse(500, "Error writing to file")))
		return
	}

	fmt.Println("Файл успешно создан и записан:", fileName)
}
