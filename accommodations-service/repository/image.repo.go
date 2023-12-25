package repository

import (
	"fmt"
	hdfs "github.com/colinmarc/hdfs/v2"
	"io"
	"log"
	"mime/multipart"
	"os"
)

const (
	hdfsRoot     = "/hdfs"
	hdfsWriteDir = "/hdfs/created/"
)

type FileStorage struct {
	client *hdfs.Client
	logger *log.Logger
}

func NewFileStorage(logger *log.Logger) *FileStorage {
	hdfsUri := os.Getenv("HDFS_URI")
	client, err := hdfs.New(hdfsUri)
	if err != nil {
		log.Println(hdfsUri)

		log.Println("CRKO SAM")
		log.Println(err)
		return nil
	}

	return &FileStorage{
		client: client,
		logger: logger,
	}
}

func (fs *FileStorage) Close() {
	fs.client.Close()
}

func (fs *FileStorage) CreateDirectories() error {
	err := fs.client.MkdirAll(hdfsWriteDir, 0644)
	if err != nil {
		fs.logger.Println(err)
		return err
	}
	return nil
}

func (fs *FileStorage) WalkDirectories() []string {
	var paths []string
	callbackFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fs.logger.Printf("Directory: %s\n", path)
			path = fmt.Sprintf("Directory: %s\n", path)
			paths = append(paths, path)
		} else {
			fs.logger.Printf("File: %s\n", path)
			path = fmt.Sprintf("File: %s\n", path)
			paths = append(paths, path)
		}
		return nil
	}
	fs.client.Walk(hdfsRoot, callbackFunc)
	return paths
}

func (fs *FileStorage) WriteFile(fileContent multipart.File, fileName string) error {
	filePath := hdfsWriteDir + fileName
	log.Println("HDFS PATH IS:", filePath)
	file, err := fs.client.Create(filePath)
	if err != nil {
		fs.logger.Println("Error in creating file on HDFS:", err)
		return err
	}
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		fs.logger.Println("Convert error:", err)
		return err
	}

	_, err = file.Write(fileBytes)
	if err != nil {
		fs.logger.Println("Error in writing file on HDFS:", err)
		return err
	}
	_ = file.Close()
	return nil
}

func (fs *FileStorage) ReadFile(fileName string) ([]byte, error) {
	filePath := hdfsWriteDir + fileName
	file, err := fs.client.Open(filePath)
	if err != nil {
		fs.logger.Println("Error in opening file for reading on HDFS:", err)
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		fs.logger.Println("Error reading file content:", err)
		return nil, err
	}

	return fileContent, nil
}
