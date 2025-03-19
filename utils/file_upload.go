package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	// Ensure the upload directory exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.MkdirAll(folder, os.ModePerm) // creates all parent directories if needed
		if err != nil {
			return "", errors.New("failed to create upload directory")
		}
	}

	// Generate unique file name
	uniqueFileName := generateUniqueFilename(file.Filename)
	filePath := filepath.Join(folder, uniqueFileName)
	// Save the file
	if err := saveFile(file, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func RemoveFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func saveFile(file *multipart.FileHeader, filePath string) error {
	src, err := file.Open()
	if err != nil {
		return errors.New("failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return errors.New("failed to create file on disk")
	}
	defer dst.Close()

	// Copy file content
	if _, err := dst.ReadFrom(src); err != nil {
		return errors.New("failed to save file content")
	}
	return nil
}

func generateUniqueFilename(originalName string) string {
	// timestamp := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	randomNumber := rand.Intn(100000)
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)] // Get file name without extension
	return fmt.Sprintf("%s_%d%s%s", timestampStr, randomNumber, name, ext)
}
