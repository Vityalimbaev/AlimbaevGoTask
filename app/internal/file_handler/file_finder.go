package file_handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func FindSmallestLargestFileNames(dirPath string) (string, string, bool) {

	fileNames := getAndFilterFilesNameInDir(dirPath)

	if len(fileNames) < 2 {
		log.Printf("warning: Not enough files found in the \"%s\"", dirPath)
		return "", "", false
	}

	min := int(^uint64(0) >> 1)
	max := -min - 1

	for _, name := range fileNames {

		num, _ := strconv.Atoi(name)

		if num < min {
			min = num
		}

		if num > max {
			max = num
		}
	}

	log.Printf("Max name: %d; Min name: %d", max, min)

	maxFilePath := fmt.Sprintf("%s/%d.log", dirPath, max)
	minFilePath := fmt.Sprintf("%s/%d.log", dirPath, min)

	return minFilePath, maxFilePath, true
}

func getAndFilterFilesNameInDir(dirPath string) []string {

	var fileNames []string

	_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			log.Printf("error: %s ", err.Error())
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".log" {

			name := strings.Split(info.Name(), ".log")[0]

			if _, err = strconv.Atoi(name); err == nil {
				fileNames = append(fileNames, name)
			}
		}

		return nil
	})

	return fileNames
}
