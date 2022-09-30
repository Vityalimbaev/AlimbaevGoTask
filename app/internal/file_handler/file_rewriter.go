package file_handler

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

type ChanMsg struct {
	offset int
	value  byte
}

type FileMutex struct {
	sync.Mutex
	file *os.File
}

func RewriteFiles(minFilePath, maxFilePath string) error {

	file1, _ := os.OpenFile(minFilePath, os.O_RDWR, os.ModeAppend)
	file2, _ := os.OpenFile(maxFilePath, os.O_RDWR, os.ModeAppend)

	stat1, _ := file1.Stat()
	stat2, _ := file2.Stat()

	fileSize1 := stat1.Size()
	fileSize2 := stat2.Size()

	var maxSize int64 = 0
	if fileSize1 > fileSize2 {
		maxSize = fileSize1
	} else {
		maxSize = fileSize2
	}

	defer func(f1 *os.File, f2 *os.File) {
		_ = f1.Close()
		_ = f2.Close()
	}(file1, file2)

	ch1 := make(chan byte)
	ch2 := make(chan byte)

	var wg sync.WaitGroup

	go w1(ch1, ch2, maxSize, file1, &wg)
	go w1(ch2, ch1, maxSize, file2, &wg)

	wg.Add(2)

	log.Println("start of rewriting")
	ch1 <- 0

	wg.Wait()

	_ = file1.Truncate(fileSize2)
	_ = file2.Truncate(fileSize1)

	log.Println("rewrite is done")

	return nil
}

func w1(in <-chan byte, out chan<- byte, maxCountOfIterations int64, file *os.File, group *sync.WaitGroup) {

	defer func() {
		close(out)
		group.Done()
	}()

	var writeByteCount int64 = 0
	writeFunc := func(msg byte) {
		if msg != 0 {
			if _, err := file.WriteAt([]byte{msg}, writeByteCount); err != nil {
				log.Println("write error : ", err)
			} else {
				writeByteCount++
			}

		}
	}

	r := bufio.NewReader(file)
	readFunc := func() byte {
		b, err := r.ReadByte()

		if err != nil && !errors.Is(err, io.EOF) {

			log.Printf("read error (%v) : %v", file.Name(), err)
		}

		return b
	}

	var countOfIterations int64 = 0

	for {
		select {
		case msg, _ := <-in:
			if countOfIterations > maxCountOfIterations {
				return
			}

			b := readFunc()

			if msg != 0 {
				writeFunc(msg)
			}

			out <- b
			countOfIterations++
		}
	}
}
