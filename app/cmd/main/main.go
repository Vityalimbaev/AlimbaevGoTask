package main

import (
	"FileHandlerTaskProject/app/config"
	"FileHandlerTaskProject/app/internal/file_handler"
)

func main() {
	conf := config.GetConfig()

	minFilePath, maxFilePath, ok := file_handler.FindSmallestLargestFileNames(conf.DataPath)

	if ok {
		file_handler.RewriteFiles(minFilePath, maxFilePath)
	}
}
